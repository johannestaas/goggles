module github.com/johannestaas/goggles

go 1.16

require (
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mdempsky/gocode v0.0.0-20200405233807-4acdcbdea79d // indirect
	github.com/motemen/gore v0.5.2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.0 // indirect
	internal/connection v0.0.1
	internal/database v0.0.1
	internal/kvstore v0.0.1
)

replace internal/connection => ./internal/connection

replace internal/kvstore => ./internal/kvstore

replace internal/database => ./internal/database
