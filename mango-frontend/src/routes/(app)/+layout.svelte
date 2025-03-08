<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import {apiBaseUrl} from '$lib/stores/config'
  import "../../app.css"
  // User profile state
  let username = "";
  let userImage = "";
  
  // Search functionality
  let searchQuery = "";
  
  // Function to handle search submission
  function handleSearch() {
    if (searchQuery.trim()) {
      goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  }
  
  // Fetch user profile data
  onMount(async () => {
    try {

      const res = await fetch(`${$apiBaseUrl}/user/profile`, {
        credentials: 'include'
      });
      
      if (res.ok) {
        const userData = await res.json();
        username = userData.username;
        userImage = 'https://ia800305.us.archive.org/31/items/discordprofilepictures/discordgrey.png';
      } else {
        // If not logged in, redirect to login page
        // goto('/login');
        username = 'para';
        userImage = 'https://ia800305.us.archive.org/31/items/discordprofilepictures/discordgrey.png';
      }
    } catch (error) {
      console.error("Failed to fetch user profile:", error);
    }
  });
</script>

<div class="min-h-screen bg-gray-900 text-white">
  <!-- Navigation Bar -->
  <header class="bg-black bg-opacity-90 fixed w-full z-50">
    <div class="container mx-auto px-4 py-3">
      <div class="flex items-center justify-between">
        <!-- Logo and Main Navigation -->
        <div class="flex items-center space-x-8">
          <!-- Logo -->
          <a href="/" class="text-red-600 font-bold text-2xl">MediaHub</a>
          
          <!-- Main Navigation Links -->
          <nav class="hidden md:flex space-x-6">
            <a href="/home" class="hover:text-gray-300 transition duration-200">Home</a>
            <a href="/manga" class="hover:text-gray-300 transition duration-200">Manga</a>
            <a href="/audiobooks" class="hover:text-gray-300 transition duration-200">Audiobooks</a>
            <a href="/anime" class="hover:text-gray-300 transition duration-200">Anime</a>
            <a href="/movies" class="hover:text-gray-300 transition duration-200">Movies</a>
          </nav>
        </div>
        
        <!-- Search and User Profile -->
        <div class="flex items-center space-x-4">
          <!-- Search Bar -->
          <form on:submit|preventDefault={handleSearch} class="relative">
            <input 
              type="text" 
              bind:value={searchQuery}
              placeholder="Search titles..." 
              class="bg-gray-800 rounded-full py-1 px-4 text-sm focus:outline-none focus:ring-1 focus:ring-red-600 w-40 md:w-64"
            />
            <button aria-label='search' type="submit" class="absolute right-3 top-1/2 transform -translate-y-1/2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
          </form>
          
          <!-- User Profile -->
          {#if username}
            <div class="flex items-center space-x-2 group relative">
              <div class="h-8 w-8 rounded-full overflow-hidden border-2 border-transparent group-hover:border-white transition duration-200">
                {#if userImage}
                  <img src={userImage} alt={username} class="h-full w-full object-cover" />
                {:else}
                  <div class="h-full w-full bg-red-600 flex items-center justify-center text-white">
                    {username.charAt(0).toUpperCase()}
                  </div>
                {/if}
              </div>
              <span class="hidden md:inline">{username}</span>
              <!-- Dropdown Menu (appears on hover) -->
              <div class="absolute hidden group-hover:block right-0 top-full mt-2 w-48 bg-black bg-opacity-90 rounded shadow-lg">
                <a href="/profile" class="block px-4 py-2 hover:bg-gray-800">Profile</a>
                <a href="/settings" class="block px-4 py-2 hover:bg-gray-800">Settings</a>
                <button 
                  class="block w-full text-left px-4 py-2 hover:bg-gray-800" 
                  on:click={async () => {
                    await fetch(`${$apiBaseUrl}/logout`, {
                      method: 'POST',
                      credentials: 'include'
                    });
                    goto('/');
                  }}
                >
                  Logout
                </button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </header>
  
  <!-- Content Area (with padding for fixed header) -->
  <main class="pt-16">
    <slot></slot>
  </main>
  
  <!-- Footer -->
  <footer class="bg-black bg-opacity-80 py-6 mt-12">
    <div class="container mx-auto px-4">
      <div class="text-center text-gray-400 text-sm">
        <p>&copy; {new Date().getFullYear()} MediaHub. All rights reserved.</p>
      </div>
    </div>
  </footer>
</div>