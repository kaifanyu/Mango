import type { LayoutLoad } from './$types';
import { goto } from '$app/navigation';

export const load: LayoutLoad = async ({ fetch }) => {
    const res = await fetch('http://localhost:8080/auth/check', { credentials: 'include' });

    if (!res.ok) {
        console.log("Auth check failed")
        goto('/'); // Redirect to login if not authenticated
        return { user: null };
    }

    const user = await res.json();
    return { user };
};
