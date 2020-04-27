package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/framework"
)

// Prefix command
func Prefix(ctx *exrouter.Context) {
	// Direct messages
	if ctx.Msg.GuildID == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You can't set a prefix in a DM.")
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
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You do not have administrator permissions to change prefix.")
		return
	}

	newPrefix := ctx.Args.After(1)

	if newPrefix == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You did not provide a new prefix.")
		return
	}

	framework.PDB.UpdateGuild(ctx.Msg.GuildID, newPrefix)
	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Set the new prefix to "+newPrefix)
}
