package constant

type UserRole string
type UserTier string

const (
    RoleAdmin     UserRole = "admin"
    RoleUser      UserRole = "user"
    RoleModerator UserRole = "moderator"
)

const (
    TierBronze    UserTier = "bronze"
    TierSilver    UserTier = "silver"
    TierGold      UserTier = "gold"
    TierDiamond   UserTier = "diamond"
    TierPlatinum  UserTier = "platinum"
)