-- This is Initialization file for DB
-- CREATE TABLE IF NOT EXISTS instance (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   client_id TEXT,
--   display_label TEXT,
--   short_label TEXT,
--   size_gb REAL,
--   num_nodes INTEGER, -- number of nodes for the instance, 1 => only master
--   api_key TEXT,
--   client_id_key TEXT,
--   client_secret TEXT,
--   note TEXT,
--   description TEXT,
--   last_updated DATE,
--   created DATE
-- );

-- CREATE TABLE IF NOT EXISTS instance_detail (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   instance_id INTEGER,
--   server_id INTEGER, -- sureSQL Node ID, which VPN server it is deployed
--   node_id INTEGER, -- rqlite node-id
--   size_mb REAL, -- actual size now
--   aplication TEXT, -- rqlite, suresql-backend, any other later on like redis, nats, etc
--   last_updated DATE,
--   last_checked DATE,
--   created DATE
-- );

-- Init DB tables for each node, settings of current node
-- DROP TABLE _settings;
CREATE TABLE IF NOT EXISTS _configs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  label TEXT, -- the name of the project
  ip TEXT,
  host TEXT,
  port TEXT,
  ssl BOOLEAN,
  dbms TEXT, -- rqlite, mysql, postgres
  mode TEXT, -- r, w, rw, b (backup)
  nodes INTEGER, -- total number of nodes in this project
  node_number INTEGER, -- this node serial number, if 1 then it's master!
  is_init_done BOOLEAN, -- database already initialized
  is_split_write BOOLEAN, -- write and read queries are separated to different nodes
  encryption_method TEXT -- none/AES/BCrypt
);



-- Init DB tables for each node, anything that is more dynamic to be put into settings
-- Use this to store the peers node as following:
-- category: nodes
-- data_type: text
-- config_key: node_name
-- text_value: node_number;hostname;ip;mode
-- example:
--   config_key: node_master
--   text_value: 1;node_master.suresql.app;120.13.217.10;r
-- example:
--   config_key: TOKEN_EXP_MINUTES
--   int_value: 60
CREATE TABLE IF NOT EXISTS _settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  category TEXT, -- grouping of configs
  data_type TEXT,  -- int/float/string/bool (which is int)
  setting_key TEXT, -- the key in Key-Value map
  text_value TEXT,
  float_value REAL,
  int_value INTEGER
);

-- Init DB for logging
CREATE TABLE IF NOT EXISTS _access_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT, 
  action_type TEXT,   -- select, insert, update, delete
  occurred TEXT DEFAULT CURRENT_TIMESTAMP,
  table_name TEXT,
  raw_query TEXT,     -- is this going to be too big?
  result TEXT,
  result_status TEXT, -- success/failed/etc
  error TEXT,         -- if there is error
  duration REAL,      -- in ms
  method TEXT,        -- API/GET/POST
  node_number	INTEGER,	  -- just in case it's merged from all nodes, this shows which is coming from
  note TEXT,
  description TEXT,
  client_ip	TEXT,	    -- from which IP
  client_browser	TEXT,	-- from which browser
  client_device	TEXT	-- from which device (mobile;xiaomi type)
);


-- User and Token, password is hashed in HashPin
CREATE TABLE IF NOT EXISTS _users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT,
  password TEXT, -- hashed
  role_name TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- token is using medaLib NewToken or encrypted. NOTE: unused for the moment, use TTLMap
CREATE TABLE IF NOT EXISTS _tokens (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id TEXT, -- if using rqlite then this is username
  token TEXT,
  refresh TEXT,
  token_expired_at TEXT,
  refresh_expired_at TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- for buckets and files
CREATE TABLE IF NOT EXISTS buckets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  label TEXT,
  short_label TEXT,
  category TEXT,
  parent INT, -- parent id which is bucket_id, this is for nesting
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  bucket_id INT,
  file_name TEXT,
  label TEXT,
  short_label TEXT,
  file_type TEXT,
  file_size FLOAT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS _acl_file (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  file_id INT,
  role_id INT,
  access_create BOOLEAN,
  access_read BOOLEAN,
  access_update BOOLEAN,
  access_delete BOOLEAN,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS _acl_bucket (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  bucket_id INT,
  role_id INT,
  access_create BOOLEAN,
  access_read BOOLEAN,
  access_update BOOLEAN,
  access_delete BOOLEAN,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- acl name like 'db admin' then add this role_id into acl_[something] like acl_file
-- TODO: add _acl_db or _acl_table for access to database in general or to scope level down to 'table'
CREATE TABLE IF NOT EXISTS _acl_role (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  label TEXT,
  short_label TEXT,
  description TEXT,
  category TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
