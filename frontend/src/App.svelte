<script lang="ts">
    import logo from "./assets/images/logo-universal.png";
    import { StartBot, FetchMembers } from "../wailsjs/go/main/App.js";
    import * as ConfigStore from "../wailsjs/go/wailsconfigstore/ConfigStore";

    let botToken: string;
    let voiceChannelID: string;
    let guildID: string;
    let obsPassword: string;
    let membersContainer: HTMLDivElement;

    async function fetchMembers(): Promise<void> {
        const membersMap = await FetchMembers();
        console.log(membersMap, typeof membersMap);
        // Build HTML string
        let html = "";
        for (const [id, username] of membersMap.entries()) {
            html += `
                <div class="member">
                    <div class="member-info">
                        <h3>${username}</h3>
                        <p>ID: ${id}</p>
                    </div>
                </div>
            `;
        }
        // Directly set innerHTML
        if (membersContainer) {
            membersContainer.innerHTML = html;
        }
    }

    async function submitBotToken(): Promise<void> {
        await ConfigStore.Set(
            "botData",
            JSON.stringify({
                token: botToken,
                voiceChannelID: voiceChannelID,
                guildID: guildID,
                obsPassword: obsPassword,
            }),
        );
        isBotSetUp = true;
        await StartBot(botToken, voiceChannelID, guildID, obsPassword);
    }
    ConfigStore.Get("botData", "null").then((response) => {
        const data = JSON.parse(response);
        botToken = data.token;
        voiceChannelID = data.voiceChannelID;
        guildID = data.guildID;
        obsPassword = data.obsPassword;
    });
    let isBotSetUp = botToken != "" && botToken != null;
</script>

<main>
    {#if isBotSetUp}
        <div class="result" id="result">Bot is set up</div>
        <button on:click={fetchMembers}>Fetch Members</button>

        <div class="members-container" bind:this={membersContainer}></div>
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
        <input
            type="text"
            name="obsPassword"
            placeholder="Enter OBS password"
            bind:value={obsPassword}
        />
        <button on:click={submitBotToken}>Greet</button>
    {/if}
</main>

<style>
</style>
