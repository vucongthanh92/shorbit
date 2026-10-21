CREATE TABLE submissions (
  id UUID PRIMARY KEY,
  owner_user_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  content JSONB NOT NULL DEFAULT '{}'::jsonb,
  share_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE share_links (
  id UUID PRIMARY KEY,
  submission_id UUID NOT NULL REFERENCES submissions(id),
  short_code VARCHAR(32) NOT NULL UNIQUE,
  status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked')),
  access_level VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (access_level IN ('public', 'authenticated', 'restricted')),
  view_count BIGINT NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_by UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_share_links_submission_id ON share_links(submission_id);

CREATE TABLE share_link_view_daily (
  share_link_id UUID NOT NULL REFERENCES share_links(id),
  viewed_date DATE NOT NULL,
  total_views BIGINT NOT NULL DEFAULT 0,
  unique_views BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (share_link_id, viewed_date)
);
