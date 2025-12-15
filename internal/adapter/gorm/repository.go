package gorm

import (
	"context"
	"log"
	"os"
	"spahtmx/internal/domain"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormRepository struct {
	DB *gorm.DB
}

type Base struct {
 ID        string     `gorm:"type:uuid;primary_key;"`
 CreatedAt time.Time  `json:"created_at"`
 UpdatedAt time.Time  `json:"updated_at"`
 DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
 b.ID = ksuid.New().String()
 return
}

type UserGorm struct {
	Base
	Username string
	Email    string
	Status   bool
}

func (u *UserGorm) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

func toDomainUsers(users []UserGorm) []domain.User {
	var domainUsers []domain.User
	for _, u := range users {
		domainUsers = append(domainUsers, u.ToDomain())
	}
	return domainUsers
}	

func FromDomain(user domain.User) *UserGorm {
	return &UserGorm{
		Base: Base{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt,
		},
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}
}

func InitDB() *gorm.DB {

	newLogger := logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,        // Don't include params in the SQL log
		Colorful:                  false,       // Disable color
	},
	)

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&UserGorm{})
	if err != nil {
		return nil
	}

	count, err := gorm.G[UserGorm](database).Count(context.Background(), "id")
	if err != nil {
		log.Fatal("count failed", err)
	}
	log.Printf("Count %d", count)
	if count == 0 {
		users := []UserGorm{
			{Username: "alice", Email: "alice@fake.com", Status: true},
			{Username: "bob", Email: "bob@fake.com", Status: false},
			{Username: "charlie", Email: "charlie@fake.com", Status: true},
		}
		err := gorm.G[[]UserGorm](database).Create(context.Background(), &users)
		if err != nil {
			log.Fatal("insert seed records failed", err)
		}

	}

	return database
}

func NewGormRepository(db *gorm.DB) *GormRepository {

	return &GormRepository{
		DB: db,
	}
}

func (g *GormRepository) GetUsers() []domain.User {

	var users []UserGorm
	g.DB.Find(&users)
	return toDomainUsers(users)
}

func (g *GormRepository) UpdateUserStatus(ctx context.Context, id string) {

	genericDB := gorm.G[UserGorm](g.DB)
	user, err := genericDB.Where("id = ?", id).Take(ctx)
	if err != nil {
		return
	}
	_, err = genericDB.Where("id = ?", id).Update(ctx, "status", !user.Status)
	if err != nil {
		return
	}
}

func (g *GormRepository) GetUserCount() string {
	return "210"
}

func (g *GormRepository) GetPageView() string {
	return "12345"
}


