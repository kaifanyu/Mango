<script lang="ts">
    import { onMount } from 'svelte';
    
    // Define TypeScript interfaces for content data
    interface ContentItem {
      id: number;
      type: string;
      title: string;
      description?: string;
      cover_image?: string;
      year?: number;
      author?: string;
    }
  
    // Content state with proper typing
    let manga: ContentItem[] = [];
    let audiobooks: ContentItem[] = [];
    let anime: ContentItem[] = [];
    let movies: ContentItem[] = [];
    let featuredContent: ContentItem | null = null;
    let isLoading = true;
    let error: string | null = null;
    // Carousel functionality
    function scrollCarousel(id: string, direction: 'left' | 'right'): void {
      const carousel = document.getElementById(id);
      if (carousel) {
        const scrollAmount = direction === 'left' ? -carousel.offsetWidth * 0.8 : carousel.offsetWidth * 0.8;
        carousel.scrollBy({ left: scrollAmount, behavior: 'smooth' });
      }
    }
    
    // Fetch content data from the backend
    onMount(async () => {
      try {
        isLoading = true;
        
        // Fetch data for each content type
        const [mangaRes, audiobooksRes, animeRes, moviesRes] = await Promise.all([
          fetch('http://localhost:8080/api/content/manga'),
          fetch('http://localhost:8080/api/content/audiobooks'),
          fetch('http://localhost:8080/api/content/anime'),
          fetch('http://localhost:8080/api/content/movies')
        ]);
        
        // Process responses
        if (mangaRes.ok) manga = await mangaRes.json();
        if (audiobooksRes.ok) audiobooks = await audiobooksRes.json();
        if (animeRes.ok) anime = await animeRes.json();
        if (moviesRes.ok) movies = await moviesRes.json();
        
        // Set a random featured content from all items
        const allContent = [...manga, ...audiobooks, ...anime, ...movies];
        if (allContent.length > 0) {
          featuredContent = allContent[Math.floor(Math.random() * allContent.length)];
        }
        console.log("THIS: ", featuredContent?.cover_image)
      } catch (err) {
        console.error("Error fetching content:", err);
        error = "Failed to load content. Please try again later.";
      } finally {
        isLoading = false;
      }
    });
  </script>
  
  <div class="bg-gray-900 text-white min-h-screen">
    <!-- Hero/Featured Content Banner -->
    {#if featuredContent}
      <div class="relative h-[50vh] md:h-[70vh] overflow-hidden">
        <div class="absolute inset-0 bg-gradient-to-r from-black to-transparent z-10"></div>
        <img 
          src={ featuredContent.cover_image || '/placeholder-hero.jpg'} 
          alt={featuredContent.title} 
          class="w-full h-full object-cover"
        />
        <div class="absolute bottom-0 left-0 p-6 md:p-12 z-20 max-w-xl">
          <h1 class="text-3xl md:text-5xl font-bold mb-3">{featuredContent.title}</h1>
          <p class="text-sm md:text-base mb-4 text-gray-300">{featuredContent.description || 'No description available'}</p>
          <div class="flex space-x-3">
            <a 
              href={`/${featuredContent.type}/${featuredContent.id}`} 
              class="bg-red-600 px-5 py-2 rounded font-bold hover:bg-red-700 transition duration-200"
            >
              Watch Now
            </a>
            <button 
              class="bg-gray-700 bg-opacity-70 px-5 py-2 rounded font-bold hover:bg-opacity-90 transition duration-200"
            >
              + My List
            </button>
          </div>
        </div>
      </div>
    {/if}
  
    <!-- Loading State -->
    {#if isLoading}
      <div class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-red-600"></div>
      </div>
    {/if}
    
    <!-- Error Message -->
    {#if error}
      <div class="text-center py-12 px-4">
        <p class="text-red-500">{error}</p>
        <button 
          on:click={() => window.location.reload()}
          class="mt-4 bg-red-600 px-4 py-2 rounded text-white hover:bg-red-700"
        >
          Retry
        </button>
      </div>
    {/if}
    
    <!-- Content Sections -->
    <div class="container mx-auto px-4 py-8 space-y-12">
      <!-- Manga Section -->
      {#if manga.length > 0}
        <section>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl md:text-2xl font-bold">Manga</h2>
            <a href="/manga" class="text-sm text-gray-400 hover:text-white">View All</a>
          </div>
          <div class="relative">
            <!-- Left Scroll Button -->
            <button 
              on:click={() => scrollCarousel('manga-carousel', 'left')}
              class="absolute left-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll left"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <!-- Carousel -->
            <div id="manga-carousel" class="flex overflow-x-auto space-x-4 pb-4 scrollbar-hide snap-x">
              {#each manga as item}
                <a href={`/manga/${item.id}`} class="snap-start flex-none w-32 sm:w-40 md:w-48 transition transform hover:scale-105">
                  <div class="relative aspect-[2/3] rounded overflow-hidden">
                    <img 
                      src={item.cover_image || '/placeholder-cover.jpg'} 
                      alt={item.title} 
                      class="w-full h-full object-cover"
                      loading="lazy"
                    />
                  </div>
                  <h3 class="mt-2 text-sm truncate">{item.title}</h3>
                </a>
              {/each}
            </div>
            
            <!-- Right Scroll Button -->
            <button 
              on:click={() => scrollCarousel('manga-carousel', 'right')}
              class="absolute right-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll right"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </section>
      {/if}
      
      <!-- Audiobooks Section -->
      {#if audiobooks.length > 0}
        <section>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl md:text-2xl font-bold">Audiobooks</h2>
            <a href="/audiobooks" class="text-sm text-gray-400 hover:text-white">View All</a>
          </div>
          <div class="relative">
            <!-- Left Scroll Button -->
            <button 
              on:click={() => scrollCarousel('audiobooks-carousel', 'left')}
              class="absolute left-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll left"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <!-- Carousel -->
            <div id="audiobooks-carousel" class="flex overflow-x-auto space-x-4 pb-4 scrollbar-hide snap-x">
              {#each audiobooks as item}
                <a href={`/audiobooks/${item.id}`} class="snap-start flex-none w-32 sm:w-40 md:w-48 transition transform hover:scale-105">
                  <div class="relative aspect-square rounded overflow-hidden">
                    <img 
                      src={item.cover_image || '/placeholder-cover.jpg'} 
                      alt={item.title} 
                      class="w-full h-full object-cover"
                      loading="lazy"
                    />
                  </div>
                  <h3 class="mt-2 text-sm truncate">{item.title}</h3>
                </a>
              {/each}
            </div>
            
            <!-- Right Scroll Button -->
            <button 
              on:click={() => scrollCarousel('audiobooks-carousel', 'right')}
              class="absolute right-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll right"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </section>
      {/if}
      
      <!-- Anime Section -->
      {#if anime.length > 0}
        <section>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl md:text-2xl font-bold">Anime</h2>
            <a href="/anime" class="text-sm text-gray-400 hover:text-white">View All</a>
          </div>
          <div class="relative">
            <!-- Left Scroll Button -->
            <button 
              on:click={() => scrollCarousel('anime-carousel', 'left')}
              class="absolute left-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll left"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <!-- Carousel -->
            <div id="anime-carousel" class="flex overflow-x-auto space-x-4 pb-4 scrollbar-hide snap-x">
              {#each anime as item}
                <a href={`/anime/${item.id}`} class="snap-start flex-none w-32 sm:w-40 md:w-48 transition transform hover:scale-105">
                  <div class="relative aspect-video rounded overflow-hidden">
                    <img 
                      src={item.cover_image || '/placeholder-cover.jpg'} 
                      alt={item.title} 
                      class="w-full h-full object-cover"
                      loading="lazy"
                    />
                  </div>
                  <h3 class="mt-2 text-sm truncate">{item.title}</h3>
                </a>
              {/each}
            </div>
            
            <!-- Right Scroll Button -->
            <button 
              on:click={() => scrollCarousel('anime-carousel', 'right')}
              class="absolute right-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll right"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </section>
      {/if}
      
      <!-- Movies Section -->
      {#if movies.length > 0}
        <section>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl md:text-2xl font-bold">Movies</h2>
            <a href="/movies" class="text-sm text-gray-400 hover:text-white">View All</a>
          </div>
          <div class="relative">
            <!-- Left Scroll Button -->
            <button 
              on:click={() => scrollCarousel('movies-carousel', 'left')}
              class="absolute left-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll left"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <!-- Carousel -->
            <div id="movies-carousel" class="flex overflow-x-auto space-x-4 pb-4 scrollbar-hide snap-x">
              {#each movies as item}
                <a href={`/movies/${item.id}`} class="snap-start flex-none w-32 sm:w-40 md:w-48 transition transform hover:scale-105">
                  <div class="relative aspect-[2/3] rounded overflow-hidden">
                    <img 
                      src={item.cover_image || '/placeholder-cover.jpg'} 
                      alt={item.title} 
                      class="w-full h-full object-cover"
                      loading="lazy"
                    />
                  </div>
                  <h3 class="mt-2 text-sm truncate">{item.title}</h3>
                </a>
              {/each}
            </div>
            
            <!-- Right Scroll Button -->
            <button 
              on:click={() => scrollCarousel('movies-carousel', 'right')}
              class="absolute right-0 top-1/2 transform -translate-y-1/2 z-10 bg-black bg-opacity-50 rounded-full p-2 hover:bg-opacity-75 focus:outline-none hidden md:block"
              aria-label="Scroll right"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </section>
      {/if}
    </div>
  </div>
  
  <style>
    /* Hide scrollbar for Chrome, Safari and Opera */
    .scrollbar-hide::-webkit-scrollbar {
      display: none;
    }
    
    /* Hide scrollbar for IE, Edge and Firefox */
    .scrollbar-hide {
      -ms-overflow-style: none;  /* IE and Edge */
      scrollbar-width: none;  /* Firefox */
    }
  </style>