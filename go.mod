module github.com/johannestaas/goggles

go 1.13

require internal/connection v0.0.1

require internal/kvstore v0.0.1

replace internal/connection => ./internal/connection

replace internal/kvstore => ./internal/kvstore
