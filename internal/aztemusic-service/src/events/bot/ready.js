const { ActivityType } = require('discord.js')
const logger = require("../../utils/logger");
const config = require("../../config");

module.exports = {
    name: "ready",
    once: true,
    async execute(client) {
        if (config.analytics) {
            const fetch = require("node-fetch");
            const crypto = require("node:crypto");

            let options = {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: `{"identifier":"${crypto.createHash("sha256").update(config.clientId).digest("hex")}"}`,
            };

            fetch("https://analytics.techy.lol/melody", options).catch(() => {});
        }

        logger.success("AzteRadio is now ready.");

        if (client.guilds.cache.size === 0) {
            logger.warn(`AzteRadio is not in any servers. Invite AzteRadio to your server using the following link: https://discord.com/api/oauth2/authorize?client_id=${config.clientId}&permissions=274914887744&scope=bot%20applications.commands`);
        } else {
            logger.info(`AzteRadio is in ${client.guilds.cache.size} ${client.guilds.cache.size === 1 ? "server" : "servers"}.`);
        }

        client.user.setPresence({ 
            activities: [{ 
                name: 'Type /help', 
                type: ActivityType.Listening
            }], 
            status: 'online' 
        });
    },
};
