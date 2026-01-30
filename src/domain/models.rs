use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    pub id: String,
    pub username: String,
    pub password: String,
    pub email: String,
    pub status: bool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PrizeList {
    pub prizes: Vec<Prize>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Prize {
    pub id: String,
    pub year: String,
    pub category: String,
    pub overall_motivation: Option<String>,
    pub laureates: Option<Vec<Laureate>>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Laureate {
    pub firstname: Option<String>,
    pub surname: Option<String>,
    pub motivation: Option<String>,
    pub share: Option<String>,
}
