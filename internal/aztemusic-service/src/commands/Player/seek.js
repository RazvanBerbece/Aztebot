const { SlashCommandBuilder } = require("@discordjs/builders");
const { EmbedBuilder } = require("discord.js");
const { Player } = require("discord-player");
const config = require("../../config");

module.exports = {
    data: new SlashCommandBuilder()
        .setName("seek")
        .setDescription("Seeks the current track to the specified position.")
        .setDMPermission(false)
        .addIntegerOption((option) => option.setName("minutes").setDescription("The amount of minutes to seek to.").setRequired(true))
        .addIntegerOption((option) => option.setName("seconds").setDescription("The amount of seconds to seek to.").setRequired(true)),
    async execute(interaction) {
        const player = Player.singleton();
        const queue = player.nodes.get(interaction.guild.id);

        const embed = new EmbedBuilder();
        embed.setColor(config.embedColour);

        if (!queue || !queue.isPlaying()) {
            embed.setDescription("There isn't currently any music playing.");
            return await interaction.reply({ embeds: [embed] });
        }

        const minutes = interaction.options.getInteger("minutes");
        const seconds = interaction.options.getInteger("seconds");

        const newPosition = minutes * 60 * 1000 + seconds * 1000;

        queue.node.seek(newPosition);

        embed.setDescription(`The current track has been seeked to **${minutes !== 0 ? `${minutes} ${minutes == 1 ? "minute" : "minutes"} and ` : ""} ${seconds} ${seconds == 1 ? "second" : "seconds"}**.`);

        return await interaction.reply({ embeds: [embed] });
    },
};
