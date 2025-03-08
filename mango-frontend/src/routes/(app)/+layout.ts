import type { LayoutLoad } from './$types';
import { goto } from '$app/navigation';
import { apiBaseUrl } from '$lib/stores/config';
import { get } from 'svelte/store'; // Import the get function

export const load: LayoutLoad = async ({ fetch }) => {
    const res = await fetch(`${get(apiBaseUrl)}/auth/check`, { credentials: 'include' });

    if (!res.ok) {
        console.log("Auth check failed")
        goto('/'); // Redirect to login if not authenticated
        return { user: null };
    }

    const user = await res.json();
    return { user };
};
