const { SlashCommandBuilder, ButtonBuilder, ActionRowBuilder } = require("@discordjs/builders");
const { EmbedBuilder, ButtonStyle } = require("discord.js");
const config = require("../../config");

module.exports = {
    data: new SlashCommandBuilder().setName("botinfo").setDescription(`hows information about the ${config.appName} bot.`),
    async execute(interaction) {
        const embed = new EmbedBuilder();
        embed.setDescription(`Melody is an open source Discord music bot that can be self-hosted to allow users to listen to music, videos, and livestreams in a voice channel together. ${config.appName} builds on top of Melody by extending some features, improving performance in some cases and containerising it to facilitate CD.`);
        embed.setColor(config.embedColour);

        const row = new ActionRowBuilder().addComponents(new ButtonBuilder().setStyle(ButtonStyle.Link).setLabel("GitHub").setURL("https://github.com/NerdyTechy/Melody"), new ButtonBuilder().setStyle(ButtonStyle.Link).setLabel("Contributors").setURL("https://github.com/NerdyTechy/Melody/graphs/contributors"), new ButtonBuilder().setStyle(ButtonStyle.Link).setLabel("Support").setURL("https://github.com/NerdyTechy/Melody/discussions"), new ButtonBuilder().setStyle(ButtonStyle.Link).setLabel("Report Issue").setURL("https://github.com/NerdyTechy/Melody/issues/new/choose"));

        return await interaction.reply({ embeds: [embed], components: [row] });
    },
};
