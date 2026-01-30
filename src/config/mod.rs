use std::env;

#[derive(Clone)]
pub struct Config {
    pub port: String,
    pub mongodb_url: String,
    pub seed_db: bool,
}

impl Config {
    pub fn load() -> Self {
        dotenvy::dotenv().ok();
        Self {
            port: env::var("PORT").unwrap_or_else(|_| "8080".to_string()),
            mongodb_url: env::var("MONGODB_URL").unwrap_or_else(|_| "mongodb://root:example@localhost:27017".to_string()),
            seed_db: env::var("SEED_DB").map(|v| v == "true").unwrap_or(false),
        }
    }
}
