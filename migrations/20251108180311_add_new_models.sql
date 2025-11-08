-- +goose Up
-- +goose StatementBegin

-- Enable citext extension for case-insensitive text
CREATE EXTENSION IF NOT EXISTS citext;

-- Update users table: add bio, password, and rename avatar to avatar_url
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS bio text,
    ADD COLUMN IF NOT EXISTS password text;

-- Rename avatar column if it exists (handle if already renamed)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='users' AND column_name='avatar') THEN
        ALTER TABLE users RENAME COLUMN avatar TO avatar_url;
    END IF;
END $$;

-- Update watchlists table: add new columns and modify structure
ALTER TABLE watchlists
    ADD COLUMN IF NOT EXISTS slug citext,
    ADD COLUMN IF NOT EXISTS cover_url text,
    ADD COLUMN IF NOT EXISTS like_count int NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS save_count int NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS item_count int NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS visibility text NOT NULL DEFAULT 'private',
    ADD COLUMN IF NOT EXISTS saved_by jsonb DEFAULT '[]';

-- Generate unique slugs for existing rows if slug was just added
UPDATE watchlists SET slug = id::text WHERE slug IS NULL;

-- Now make slug NOT NULL and UNIQUE
ALTER TABLE watchlists ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS watchlists_slug_key ON watchlists(slug);

-- Add check constraint for visibility
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'watchlists_visibility_check'
    ) THEN
        ALTER TABLE watchlists 
            ADD CONSTRAINT watchlists_visibility_check 
            CHECK (visibility IN ('public', 'private', 'unlisted'));
    END IF;
END $$;

-- Drop old is_public column if it exists
ALTER TABLE watchlists DROP COLUMN IF EXISTS is_public;

-- Create movies table first (needed for foreign key)
CREATE TABLE IF NOT EXISTS movies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tmdb_id int UNIQUE NOT NULL,
    title text NOT NULL,
    year int,
    poster_url text,
    backdrop_url text,
    genres jsonb,
    runtime int,
    release_date timestamptz,
    metadata jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX IF NOT EXISTS idx_movies_year ON movies(year);
CREATE INDEX IF NOT EXISTS idx_movies_deleted_at ON movies(deleted_at);

-- Update watchlist_items table structure
-- First, handle the note/notes column rename
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='watchlist_items' AND column_name='notes') THEN
        ALTER TABLE watchlist_items RENAME COLUMN notes TO note;
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                      WHERE table_name='watchlist_items' AND column_name='note') THEN
        ALTER TABLE watchlist_items ADD COLUMN note text;
    END IF;
END $$;

-- Add movie_id if it doesn't exist
ALTER TABLE watchlist_items
    ADD COLUMN IF NOT EXISTS movie_id uuid;

-- Rename created_at to added_at if needed
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='watchlist_items' AND column_name='created_at') THEN
        ALTER TABLE watchlist_items RENAME COLUMN created_at TO added_at;
    END IF;
END $$;

-- Drop old columns from watchlist_items
ALTER TABLE watchlist_items 
    DROP COLUMN IF EXISTS tmdb_id,
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS poster_path,
    DROP COLUMN IF EXISTS release_date,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS deleted_at;

-- Create index and foreign key for movie_id
CREATE INDEX IF NOT EXISTS idx_watchlist_items_movie_id ON watchlist_items(movie_id);
ALTER TABLE watchlist_items
    DROP CONSTRAINT IF EXISTS watchlist_items_movie_id_fkey,
    ADD CONSTRAINT watchlist_items_movie_id_fkey 
        FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE;

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type text NOT NULL CHECK (type IN ('like', 'follow')),
    actor_id uuid NOT NULL,
    entity_id uuid NOT NULL,
    is_read boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Create activities table
CREATE TABLE IF NOT EXISTS activities (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type text NOT NULL CHECK (type IN ('like', 'follow', 'create_list', 'add_item')),
    subject_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_subject_id ON activities(subject_id);
CREATE INDEX IF NOT EXISTS idx_activities_created_at ON activities(created_at);

-- Create saves table (for saving watchlists)
CREATE TABLE IF NOT EXISTS saves (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    watchlist_id uuid NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, watchlist_id)
);

CREATE INDEX IF NOT EXISTS idx_saves_watchlist_id ON saves(watchlist_id);
CREATE INDEX IF NOT EXISTS idx_saves_created_at ON saves(created_at);

-- Create follows table (for following users)
CREATE TABLE IF NOT EXISTS follows (
    follower_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, followee_id)
);

CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_followee_id ON follows(followee_id);
CREATE INDEX IF NOT EXISTS idx_follows_created_at ON follows(created_at);

-- Update likes table to use composite primary key
-- First, check if likes table needs restructuring
DO $$
BEGIN
    -- Check if the old likes table has an 'id' column
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='likes' AND column_name='id'
    ) THEN
        -- Create new likes table with composite key
        CREATE TABLE likes_new (
            user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            watchlist_id uuid NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
            created_at timestamptz NOT NULL DEFAULT now(),
            PRIMARY KEY (user_id, watchlist_id)
        );
        
        -- Copy data from old likes table
        INSERT INTO likes_new (user_id, watchlist_id, created_at)
        SELECT user_id, watchlist_id, created_at 
        FROM likes
        ON CONFLICT DO NOTHING;
        
        -- Drop old table and rename new one
        DROP TABLE likes CASCADE;
        ALTER TABLE likes_new RENAME TO likes;
        
        -- Create indexes
        CREATE INDEX idx_likes_user_id ON likes(user_id);
        CREATE INDEX idx_likes_watchlist_id ON likes(watchlist_id);
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop new tables
DROP TABLE IF EXISTS follows CASCADE;
DROP TABLE IF EXISTS saves CASCADE;
DROP TABLE IF EXISTS activities CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;

-- Restore likes table with id if needed
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='likes' AND column_name='id'
    ) THEN
        CREATE TABLE likes_old (
            id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
            created_at timestamptz NOT NULL DEFAULT now(),
            user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            watchlist_id uuid NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
            UNIQUE(user_id, watchlist_id)
        );
        
        INSERT INTO likes_old (user_id, watchlist_id, created_at)
        SELECT user_id, watchlist_id, created_at FROM likes;
        
        DROP TABLE likes CASCADE;
        ALTER TABLE likes_old RENAME TO likes;
    END IF;
END $$;

-- Restore watchlist_items structure
ALTER TABLE watchlist_items
    DROP CONSTRAINT IF EXISTS watchlist_items_movie_id_fkey,
    DROP COLUMN IF EXISTS movie_id,
    DROP COLUMN IF EXISTS note,
    ADD COLUMN IF NOT EXISTS tmdb_id bigint,
    ADD COLUMN IF NOT EXISTS title text,
    ADD COLUMN IF NOT EXISTS poster_path text,
    ADD COLUMN IF NOT EXISTS release_date text,
    ADD COLUMN IF NOT EXISTS notes text,
    ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

-- Rename back
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='watchlist_items' AND column_name='added_at') THEN
        ALTER TABLE watchlist_items RENAME COLUMN added_at TO created_at;
    END IF;
END $$;

-- Drop movies table
DROP TABLE IF EXISTS movies CASCADE;

-- Restore watchlists columns
ALTER TABLE watchlists
    DROP CONSTRAINT IF EXISTS watchlists_visibility_check,
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS cover_url,
    DROP COLUMN IF EXISTS like_count,
    DROP COLUMN IF EXISTS save_count,
    DROP COLUMN IF EXISTS item_count,
    DROP COLUMN IF EXISTS visibility,
    DROP COLUMN IF EXISTS saved_by,
    ADD COLUMN IF NOT EXISTS is_public boolean NOT NULL DEFAULT true;

-- Restore users columns
ALTER TABLE users
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS password;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='users' AND column_name='avatar_url') THEN
        ALTER TABLE users RENAME COLUMN avatar_url TO avatar;
    END IF;
END $$;

-- +goose StatementEnd
