package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/framework"
)

// Prefix command
func Prefix(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)

	// Direct messages
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":information_source: You can't set a prefix in a DM.")
		return
	}

	admin := false

	// Fetch guild
	guild, _ := ctx.Ses.Guild(ctx.Msg.GuildID)

	// Check owner
	if ctx.Msg.Author.ID == guild.OwnerID {
		admin = true
	}

	// Check admin
	if !admin {
		for _, role := range ctx.Msg.Member.Roles {
			r, _ := ctx.Ses.State.Role(ctx.Msg.GuildID, role)
			if r.Permissions&discordgo.PermissionAdministrator > 0 {
				admin = true
				break
			}
		}
	}

	if !admin {
		ctx.Reply(":information_source: You do not have administrator permissions to change prefix.")
		return
	}

	newPrefix := ctx.Args.After(1)

	if newPrefix == "" {
		ctx.Reply(fmt.Sprintf("Usage: `%sprefix <prefix>`", prefix))
		return
	}

	framework.PDB.UpdateGuild(ctx.Msg.GuildID, newPrefix)
	ctx.Reply(":information_source: Set the new prefix to " + newPrefix)
}
