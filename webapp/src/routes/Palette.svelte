<script lang="ts">
    const sizes = [ 1, 3, 5, 7, 10 ]
</script>

<div class="palette">
    <div class="blacks">
        {#each sizes as size}
            <button title="Black brush of {size} pixels"><span class="size{size}"></span></button>
        {/each}
    </div>
    <div class="whites">
        {#each sizes as size}
            <button title="White brush of {size} pixels"><span class="size{size}"></span></button>
        {/each}
    </div>
</div>

<style>
    .palette {
        display: flex;
        flex-direction: row;
        gap: 0.5rem;
    }
    .palette .blacks, .palette .whites {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .palette button {
        flex-shrink: 0;
        padding: 2px;
        width: 2rem;
        height: 2rem;
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .palette button span {
        display: block;
        border-radius: 50%;
    }
    .palette button span.size1 {
        width: 1px;
        height: 1px;
    }
    .palette button span.size3 {
        width: 3px;
        height: 3px;
    }
    .palette button span.size5 {
        width: 5px;
        height: 5px;
    }
    .palette button span.size7 {
        width: 7px;
        height: 7px;
    }
    .palette button span.size10 {
        width: 10px;
        height: 10px;
    }
    .palette .blacks button {
        color: black;
        background: white;
    }
    .palette .whites button {
        color: white;
        background: black;
    }
    .palette .blacks button span {
        background: black;
    }
    .palette .whites button span {
        background: white;
    }

    @media (max-width: 800px) {
        .palette {
            flex-direction: column;
        }
        .palette .blacks, .palette .whites {
            flex-direction: row;
        }
    }
</style>
