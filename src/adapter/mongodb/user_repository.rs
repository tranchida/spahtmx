use async_trait::async_trait;
use mongodb::{Database, bson::{doc, oid::ObjectId}, options::FindOptions};
use futures_util::stream::StreamExt;
use crate::domain::models::User;
use crate::domain::repositories::UserRepository;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct UserMongo {
    #[serde(rename = "_id")]
    pub id: ObjectId,
    pub username: String,
    pub password: String,
    pub email: String,
    pub status: bool,
}

impl From<UserMongo> for User {
    fn from(u: UserMongo) -> Self {
        User {
            id: u.id.to_hex(),
            username: u.username,
            password: u.password,
            email: u.email,
            status: u.status,
        }
    }
}

pub struct UserMongoRepository {
    pub db: Database,
}

#[async_trait]
impl UserRepository for UserMongoRepository {
    async fn get_users(&self) -> Result<Vec<User>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<UserMongo>("users");
        let find_options = FindOptions::builder().limit(10).build();
        let mut cursor = collection.find(doc! {}).with_options(find_options).await?;
        let mut users = Vec::new();
        while let Some(user) = cursor.next().await {
            users.push(user?.into());
        }
        Ok(users)
    }

    async fn get_user(&self, id: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>> {
        let obj_id = ObjectId::parse_str(id)?;
        let collection = self.db.collection::<UserMongo>("users");
        let user = collection.find_one(doc! { "_id": obj_id }).await?;
        user.map(Into::into).ok_or_else(|| "User not found".into())
    }

    async fn get_by_username(&self, username: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<UserMongo>("users");
        let user = collection.find_one(doc! { "username": username }).await?;
        user.map(Into::into).ok_or_else(|| "User not found".into())
    }

    async fn create_user(&self, user: User) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<UserMongo>("users");
        let user_mongo = UserMongo {
            id: if user.id.is_empty() { ObjectId::new() } else { ObjectId::parse_str(&user.id)? },
            username: user.username,
            password: user.password,
            email: user.email,
            status: user.status,
        };
        collection.insert_one(user_mongo).await?;
        Ok(())
    }

    async fn update_user(&self, user: User) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<UserMongo>("users");
        let id = ObjectId::parse_str(&user.id)?;
        let user_mongo = UserMongo {
            id,
            username: user.username,
            password: user.password,
            email: user.email,
            status: user.status,
        };
        collection.replace_one(doc! { "_id": id }, user_mongo).await?;
        Ok(())
    }
}
