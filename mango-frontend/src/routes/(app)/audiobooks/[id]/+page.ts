import type { PageLoad } from "./$types";
import { apiBaseUrl } from "$lib/stores/config";
import { get } from 'svelte/store'; // Import the get function

export const load: PageLoad = async ({params, fetch}) => {
    const id = params.id;
    console.log("ID: ", id)
    const response = await fetch(`${get(apiBaseUrl)}/api/content/audiobooks/${id}`);
    console.log("THIS RESPONE: ", response)
    if (!response.ok){
        return {audiobook: null, error: 'Failed to load audiobook'};
    }

    const audiobook = await response.json();
    return {audiobook};
}
