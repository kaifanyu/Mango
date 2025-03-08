import type { PageLoad } from "./$types";

export const load: PageLoad = async ({params, fetch}) => {
    const id = params.id;
    console.log("ID: ", id)
    const response = await fetch(`http://localhost:8080/api/content/audiobooks/${id}`);

    if (!response.ok){
        return {audiobook: null, error: 'Failed to load audiobook'};
    }

    const audiobook = await response.json();
    return {audiobook};
}
