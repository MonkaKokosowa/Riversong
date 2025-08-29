<script lang="ts">
    import logo from "./assets/images/logo-universal.png";
    import { StartBot } from "../wailsjs/go/main/App.js";
    import * as ConfigStore from "../wailsjs/go/wailsconfigstore/ConfigStore";

    let botToken: string;
    let voiceChannelID: string;
    let guildID: string;
    async function submitBotToken(): Promise<void> {
        await ConfigStore.Set(
            "botData",
            JSON.stringify({
                token: botToken,
                voiceChannelID: voiceChannelID,
                guildID: guildID,
            }),
        );
        isBotSetUp = true;
        await StartBot(botToken, voiceChannelID, guildID);
    }
    ConfigStore.Get("botData", "null").then((response) => {
        const data = JSON.parse(response);
        botToken = data.token;
        voiceChannelID = data.voiceChannelID;
        guildID = data.guildID;
    });
    let isBotSetUp = botToken != "" && botToken != null;
</script>

<main>
    {#if isBotSetUp}
        <div class="result" id="result">Bot is set up</div>
    {:else}
        <div class="result" id="result">Bot is not set up</div>
        <input
            type="text"
            name="botToken"
            placeholder="Enter bot token"
            bind:value={botToken}
        />
        <input
            type="text"
            name="voiceChannelID"
            placeholder="Enter voice channel ID"
            bind:value={voiceChannelID}
        />
        <input
            type="text"
            name="guildID"
            placeholder="Enter guild ID"
            bind:value={guildID}
        />
        <button on:click={submitBotToken}>Greet</button>
    {/if}
</main>

<style>
</style>
