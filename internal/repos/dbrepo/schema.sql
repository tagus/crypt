PRAGMA encoding = "UTF-8";

CREATE TABLE IF NOT EXISTS crypts (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/******************************************************************************/

CREATE TABLE IF NOT EXISTS credentials (
  id TEXT PRIMARY KEY,
  current_version INTEGER NOT NULL,
  latest_version INTEGER NOT NULL,
  tags TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  accessed_at TIMESTAMP,
  accessed_count INTEGER NOT NULL,
  crypt_id TEXT NOT NULL
);

/******************************************************************************/


CREATE TABLE IF NOT EXISTS credential_versions (
  credential_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  service TEXT NOT NULL,
  domains TEXT NOT NULL,
  email TEXT,
  username TEXT,
  s_password BLOB NOT NULL,
  description TEXT NOT NULL,
  s_details BLOB NOT NULL,
  PRIMARY KEY (credential_id, version)
);
