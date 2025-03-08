<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';
    
    // Get audiobook ID from URL params
    const audiobookId = $page.params.id;
    let fileExtension = "mp3";

    // State for audiobook data
    let audiobook: {
      id: number;
      title: string;
      cover_image: string;
    } | null = null;
    
    // Audio player state
    let audioElement: HTMLAudioElement | null = null;
    let isPlaying = false;
    let currentTime = 0;
    let duration = 0;
    let volume = 0.8;
    let playbackRate = 1;
    let progressInterval: number | null = null;
    
    // UI state
    let isLoading = true;
    let isAudioLoading = true;
    let loadingStage = "Not started"; // New logging state for debug
    let error: string | null = null;
    let showVolumeControl = false;
    let showPlaybackSpeedControl = false;
    
    // Progress tracking
    let userProgress = 0;
    let isSavingProgress = false;
    let lastSavedProgress = 0;
    
    // Debug logging with timestamps
    function logDebug(message: string, data?: any): void {
      const timestamp = new Date().toISOString();
      const logPrefix = `[AudioPlayer ${timestamp}]`;
      
      if (data) {
        console.log(`${logPrefix} ${message}`, data);
      } else {
        console.log(`${logPrefix} ${message}`);
      }
      
      // Update loading stage for UI feedback
      loadingStage = message;
    }
    
    // Format time in mm:ss or hh:mm:ss
    function formatTime(seconds: number): string {
      if (isNaN(seconds)) return '00:00';
      
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      const remainingSeconds = Math.floor(seconds % 60);
      
      if (hours > 0) {
        return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
      }
      
      return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
    }
    
    // Toggle play/pause
    function togglePlay(): void {
      if (!audioElement) {
        logDebug("Toggle play called but audio element not available");
        return;
      }
      
      if (isPlaying) {
        logDebug("Pausing audio");
        audioElement.pause();
      } else {
        logDebug("Playing audio");
        // Create a promise to track play success/failure
        const playPromise = audioElement.play();
        
        if (playPromise !== undefined) {
          playPromise
            .then(() => {
              logDebug("Play started successfully");
            })
            .catch(error => {
              logDebug("Play failed", error);
              // Handle autoplay restrictions or other errors
              if (error.name === "NotAllowedError") {
                logDebug("Autoplay not allowed by browser");
              }
            });
        }
      }
      
      isPlaying = !isPlaying;
    }
    
    // Set new time position when user drags the progress bar
    function handleSeek(event: Event): void {
      if (!audioElement) {
        logDebug("Seek attempted but audio element not available");
        return;
      }
      
      const input = event.target as HTMLInputElement;
      const newTime = parseFloat(input.value);
      logDebug(`Seeking to ${newTime} seconds`);
      
      audioElement.currentTime = newTime;
      currentTime = newTime;
      
      // Save progress to server if significant change
      if (Math.abs(newTime - lastSavedProgress) > 10) {
        saveProgress(newTime);
      }
    }
    
    // Change playback speed
    function setPlaybackRate(rate: number): void {
      if (!audioElement) {
        logDebug(`Playback rate change attempted but audio element not available`);
        return;
      }
      
      logDebug(`Setting playback rate to ${rate}x`);
      playbackRate = rate;
      audioElement.playbackRate = rate;
      showPlaybackSpeedControl = false;
    }
    
    // Change volume
    function handleVolumeChange(event: Event): void {
      if (!audioElement) {
        logDebug("Volume change attempted but audio element not available");
        return;
      }
      
      const input = event.target as HTMLInputElement;
      volume = parseFloat(input.value);
      logDebug(`Setting volume to ${volume}`);
      audioElement.volume = volume;
    }
    
    // Save listening progress to server
    async function saveProgress(progress: number): Promise<void> {
      if (isSavingProgress || !audiobookId) {
        logDebug("Save progress skipped - already in progress or missing ID");
        return;
      }
      
      try {
        isSavingProgress = true;
        logDebug(`Saving progress at ${progress} seconds`);
        
        const response = await fetch('http://localhost:8080/api/progress', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            video_id: audiobookId,
            progress: Math.floor(progress),
          }),
          credentials: 'include',
        });
        
        if (response.ok) {
          lastSavedProgress = progress;
          logDebug(`Progress saved successfully`);
        } else {
          logDebug(`Failed to save progress: ${response.status} ${response.statusText}`);
          console.error('Failed to save progress');
        }
      } catch (err) {
        logDebug(`Error saving progress`, err);
        console.error('Error saving progress:', err);
      } finally {
        isSavingProgress = false;
      }
    }
    
    // Fetch audiobook data and user progress
    async function fetchData(): Promise<void> {
      try {
        isLoading = true;
        error = null;
        logDebug("Starting data fetch");
        
        // Fetch audiobook details
        logDebug(`Fetching audiobook data for ID: ${audiobookId}`);
        const audiobookResponse = await fetch(`http://localhost:8080/api/content/audiobooks/${audiobookId}`);
        
        if (!audiobookResponse.ok) {
          const errorText = await audiobookResponse.text();
          logDebug(`Audiobook fetch failed: ${audiobookResponse.status}`, errorText);
          throw new Error('Failed to load audiobook data');
        }
        
        audiobook = await audiobookResponse.json();
        logDebug("Audiobook data loaded successfully", audiobook);
        
        // Fetch user progress
        logDebug("Fetching user progress");
        const progressResponse = await fetch(`http://localhost:8080/api/progress/${audiobookId}`, {
          credentials: 'include',
        });
        
        if (progressResponse.ok) {
          const progressData = await progressResponse.json();
          userProgress = progressData.progress || 0;
          logDebug(`User progress loaded: ${userProgress} seconds`);
        } else {
          logDebug(`Progress fetch failed or no progress: ${progressResponse.status}`);
          userProgress = 0;
        }
      } catch (err) {
        logDebug("Error fetching data", err);
        console.error('Error fetching data:', err);
        error = 'Failed to load audiobook. Please try again later.';
      } finally {
        isLoading = false;
        logDebug("Data fetch completed");
      }
    }
    
    // Update progress periodically while playing
    function startProgressTracking(): void {
      if (progressInterval) {
        clearInterval(progressInterval);
        logDebug("Cleared existing progress tracking interval");
      }
      
      logDebug("Starting progress tracking interval");
      progressInterval = setInterval(() => {
        if (!audioElement || !isPlaying) return;
        
        currentTime = audioElement.currentTime;
        
        // Save progress every 30 seconds if playing
        if (Math.abs(currentTime - lastSavedProgress) > 30) {
          logDebug(`Auto-saving progress after 30s interval`);
          saveProgress(currentTime);
        }
      }, 1000) as unknown as number;
    }
    
    // Retry loading audio with optimized settings
    function retryLoadAudio(): void {
      if (!audioElement) {
        logDebug("Retry attempted but audio element not available");
        return;
      }
      
      logDebug("Retrying audio load");
      isAudioLoading = true;
      
      // Try to reset any error states
      audioElement.removeAttribute('src');
      void audioElement.load();
      
      // Small timeout to ensure browser has time to reset state
      setTimeout(() => {
        logDebug("Setting source after reset");
        audioElement!.src = `http://localhost:8080/static/content/audiobooks/${audiobookId}.${fileExtension}`;
        audioElement!.load();
      }, 100);
    }
    
    // Set up audio element and event listeners with more comprehensive event handling
    function setupAudio(): void {
      if (!audioElement || !audiobook) {
        logDebug("Setup audio called but elements not available");
        return;
      }
      
      logDebug("Setting up audio element");
      
      // Set initial state
      audioElement.volume = volume;
      audioElement.playbackRate = playbackRate;
      
      // Set preload attribute to force metadata loading
      audioElement.preload = "metadata";
      
      // Helper function to remove all listeners when needed
      const removeAllListeners = () => {
        const events = ['play', 'pause', 'ended', 'loadstart', 'loadeddata', 
                        'canplay', 'canplaythrough', 'error', 'stalled', 
                        'waiting', 'loadedmetadata', 'timeupdate', 'progress'];
                        
        events.forEach(event => {
          audioElement!.removeEventListener(event, () => {});
        });
        
        logDebug("Removed all event listeners");
      };
      
      // Playback state listeners
      audioElement.addEventListener('play', () => { 
        isPlaying = true; 
        logDebug("Event: play");
      });
      
      audioElement.addEventListener('pause', () => { 
        isPlaying = false; 
        logDebug("Event: pause");
      });
      
      audioElement.addEventListener('ended', () => { 
        isPlaying = false; 
        logDebug("Event: ended");
      });
      
      // Loading state listeners
      audioElement.addEventListener('loadstart', () => { 
        isAudioLoading = true; 
        logDebug("Event: loadstart - Audio load initiated");
      });
      
      audioElement.addEventListener('loadeddata', () => {
        logDebug("Event: loadeddata - Audio data loaded");
      });
      
      audioElement.addEventListener('canplay', () => {
        logDebug("Event: canplay - Audio can start playing");
        // Don't set isAudioLoading to false yet, wait for canplaythrough
      });
      
      audioElement.addEventListener('canplaythrough', () => {
        isAudioLoading = false;
        logDebug("Event: canplaythrough - Audio can play without buffering");
      });
      
      // Error handling
      audioElement.addEventListener('error', (e) => { 
        const error = audioElement!.error;
        isAudioLoading = false;
        
        // Detailed error logging
        if (error) {
          switch(error.code) {
            case MediaError.MEDIA_ERR_ABORTED:
              logDebug("Audio loading error: Playback aborted by user");
              break;
            case MediaError.MEDIA_ERR_NETWORK:
              logDebug("Audio loading error: Network error");
              break;
            case MediaError.MEDIA_ERR_DECODE:
              logDebug("Audio loading error: Media decode error");
              break;
            case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
              logDebug("Audio loading error: Format not supported");
              break;
            default:
              logDebug("Audio loading error: Unknown error", error);
          }
        } else {
          logDebug("Audio loading error: No error details available");
        }
        
      });
      
      audioElement.addEventListener('stalled', () => {
        logDebug("Event: stalled - Audio download stalled");
      });
      
      audioElement.addEventListener('waiting', () => {
        logDebug("Event: waiting - Audio playback waiting for more data");
      });
      
      // Metadata and timing listeners
      audioElement.addEventListener('loadedmetadata', () => {
        duration = audioElement?.duration || 0;
        logDebug(`Event: loadedmetadata - Duration: ${duration}`);
        
        if (isNaN(duration) || !isFinite(duration)) {
          logDebug("Warning: Invalid duration received");
        } else if (duration === Infinity) {
          logDebug("Warning: Duration is Infinity, possibly a stream");
        } else if (duration <= 0) {
          logDebug("Warning: Duration is zero or negative");
        } else {
          logDebug(`Duration loaded: ${formatTime(duration)}`);
        }
        
        // Resume from user's saved progress if we have valid duration
        if (userProgress > 0 && duration > 0 && userProgress < duration) {
          logDebug(`Resuming from saved position: ${userProgress} seconds`);
          audioElement!.currentTime = userProgress;
          currentTime = userProgress;
        }
      });
      
      audioElement.addEventListener('timeupdate', () => {
        // Update current time (more frequent than our interval)
        currentTime = audioElement!.currentTime;
      });
      
      // Buffer progress monitoring
      audioElement.addEventListener('progress', () => {
        // Log buffer state
        if (audioElement!.buffered.length > 0) {
          const bufferedEnd = audioElement!.buffered.end(audioElement!.buffered.length - 1);
          const bufferedPercent = (bufferedEnd / duration) * 100;
          logDebug(`Buffer progress: ${bufferedPercent.toFixed(1)}% (${formatTime(bufferedEnd)} / ${formatTime(duration)})`);
        }
      });
      
      logDebug("Starting progress tracking");
      startProgressTracking();
      
      // Force a load event to start the process
      logDebug("Forcing initial load");
      audioElement.load();
    }
    
    // Handle audio loading timeout with multiple stages
    let loadingTimeout: number | null = null;
    let extendedLoadingTimeout: number | null = null;
    
    function setupLoadingTimeout(): void {
      // Clear any existing timeouts
      if (loadingTimeout) {
        clearTimeout(loadingTimeout);
        logDebug("Cleared existing loading timeout");
      }
      
      if (extendedLoadingTimeout) {
        clearTimeout(extendedLoadingTimeout);
        logDebug("Cleared existing extended loading timeout");
      }
      
      // First timeout - initial loading warning after 10s
      loadingTimeout = setTimeout(() => {
        if (isAudioLoading) {
          logDebug("Loading taking longer than expected (10s)");
          loadingStage = "Loading taking longer than expected. Still trying...";
        }
      }, 10000) as unknown as number;
      
      // Second timeout - actual error after 40s
      extendedLoadingTimeout = setTimeout(() => {
        if (isAudioLoading && !error) {
          logDebug("Loading timeout exceeded (40s)");
          error = 'Audio is taking too long to load. Please check your connection and try again.';
          isAudioLoading = false;
        }
      }, 40000) as unknown as number;
    }
    
    // Try to fetch file size before loading (optional optimization)
    async function checkAudioFileSize(): Promise<void> {
      const formats = ["mp3", "m4b"];
      
      for (const format of formats) {
        try {
          logDebug("Checking audio file size with HEAD request");
          const response = await fetch(`http://localhost:8080/static/content/audiobooks/${audiobookId}.${format}`, {
            method: 'HEAD'
          });
          
          if (response.ok) {
            const contentLength = response.headers.get('content-length');
            if (contentLength) {
              const sizeMB = parseInt(contentLength) / (1024 * 1024);
              fileExtension = format;
              logDebug(`Audio file size: ${sizeMB.toFixed(2)} MB`);
            } else {
              logDebug("Content-Length header not available");
            }
          } else {
            logDebug(`HEAD request failed: ${response.status}`);
          }
        } catch (err) {
          logDebug("Error checking file size", err);
        }
      }

    }
    
    // Lifecycle hooks
    onMount(async () => {
      logDebug("Component mounted");
      await fetchData();
      
      if (audiobook) {
        logDebug(`Setting page title: ${audiobook.title}`);
        document.title = `${audiobook.title} | Audiobook Player`;
        
        // Optional: Check file size before loading
        await checkAudioFileSize();
        
        // Get audio element reference
        audioElement = document.getElementById('audio-player') as HTMLAudioElement;
        
        if (audioElement) {
          logDebug("Audio element found, setting up");
          setupAudio();
          setupLoadingTimeout();
        } else {
          logDebug("Audio element not found in DOM");
          error = 'Failed to initialize audio player';
        }
      }
    });
    
    onDestroy(() => {
      logDebug("Component being destroyed");
      
      // Clear all timers
      if (progressInterval) {
        clearInterval(progressInterval);
        logDebug("Progress interval cleared");
      }
      
      if (loadingTimeout) {
        clearTimeout(loadingTimeout);
        logDebug("Loading timeout cleared");
      }
      
      if (extendedLoadingTimeout) {
        clearTimeout(extendedLoadingTimeout);
        logDebug("Extended loading timeout cleared");
      }
      
      // Save final progress when leaving the page
      if (currentTime > 0 && Math.abs(currentTime - lastSavedProgress) > 5) {
        logDebug("Saving final progress before unmount");
        saveProgress(currentTime);
      }
      
      // Clean up audio element
      if (audioElement) {
        logDebug("Removing audio element source");
        audioElement.pause();
        audioElement.removeAttribute('src');
        audioElement.load();
      }
    });
  </script>
  
  <div class="bg-gray-900 text-white min-h-screen">
    <!-- Loading State for initial data -->
    {#if isLoading}
      <div class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-red-600"></div>
      </div>
    {:else if error}
      <!-- Error Message -->
      <div class="text-center py-12 px-4">
        <p class="text-red-500">{error}</p>
        <button 
          on:click={() => {
            error = null;
            fetchData();
            if (audioElement) retryLoadAudio();
          }}
          class="mt-4 bg-red-600 px-4 py-2 rounded text-white hover:bg-red-700"
        >
          Retry
        </button>
      </div>
    {:else if audiobook}
      <!-- Audiobook Player -->
      <div class="container mx-auto px-4 py-8">
        <!-- Back button -->
        <button 
          on:click={() => history.back()}
          class="mb-8 flex items-center text-gray-400 hover:text-white transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <!-- Main content -->
        <div class="md:flex">
          <!-- Cover art -->
          <div class="md:w-1/3 mb-8 md:mb-0 md:pr-8">
            <div class="rounded-lg overflow-hidden shadow-lg bg-gray-800 aspect-square">
              <img 
                src={audiobook.cover_image || '/placeholder-cover.jpg'} 
                alt={audiobook.title} 
                class="w-full h-full object-cover"
              />
            </div>
          </div>
          
          <!-- Title and player -->
          <div class="md:w-2/3">
            <h1 class="text-3xl md:text-4xl font-bold mb-6">{audiobook.title}</h1>
            
            <!-- Actual audio element (hidden) -->
            <audio 
              id="audio-player" 
              src={`http://localhost:8080/static/content/audiobooks/${audiobookId}.${fileExtension}`}
              preload="metadata"
            ></audio>
            
            <!-- Debug info (toggle in production) -->
            <!-- <div class="mb-4 p-2 bg-gray-800 rounded text-xs text-gray-400">
              <details>
                <summary class="cursor-pointer">Debug Info</summary>
                <div class="mt-2 pl-2 border-l border-gray-700">
                  <p>Loading Stage: {loadingStage}</p>
                  <p>Audio Loading: {isAudioLoading ? 'Yes' : 'No'}</p>
                  <p>Duration: {formatTime(duration)}</p>
                  <p>Is Playing: {isPlaying ? 'Yes' : 'No'}</p>
                  <p>Current Time: {formatTime(currentTime)}</p>
                  <p>User Progress: {formatTime(userProgress)}</p>
                </div>
              </details>
            </div>
             -->
  
            <!-- Audio loading state -->
            {#if isAudioLoading}
              <div class="bg-gray-800 rounded-xl p-6 shadow-lg">
                <div class="flex flex-col items-center justify-center py-8">
                  <div class="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-red-600 mb-4"></div>
                  <p class="text-gray-300">Loading audio file...</p>
                  <p class="text-gray-500 text-sm mt-2">{loadingStage}</p>
                </div>
              </div>
            {:else}
              <!-- Custom player UI - only show when audio is loaded -->
              <div class="bg-gray-800 rounded-xl p-6 shadow-lg">
                <!-- Progress bar -->
                <div class="mb-4">
                  <input 
                    type="range" 
                    min="0" 
                    max={duration || 100} 
                    value={currentTime} 
                    on:input={handleSeek}
                    class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
                  />
                  <div class="flex justify-between text-sm text-gray-400 mt-1">
                    <span>{formatTime(currentTime)}</span>
                    <span>{formatTime(duration)}</span>
                  </div>
                </div>
                
                <!-- Controls -->
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-4">
                    <!-- Skip backward 15s -->
                    <button 
                      class="text-white hover:text-red-500 transition"
                      on:click={() => {
                        if (audioElement) {
                          audioElement.currentTime = Math.max(0, audioElement.currentTime - 15);
                          currentTime = audioElement.currentTime;
                          logDebug("Skipped backward 15s");
                        }
                      }}
                      aria-label="Skip backward 15 seconds"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 5V3L8 7l4 4V9c3.31 0 6 2.69 6 6s-2.69 6-6 6-6-2.69-6-6H4c0 4.42 3.58 8 8 8s8-3.58 8-8-3.58-8-8-8z"/>
                        <text x="11" y="15" font-size="6" font-weight="bold" font-family="sans-serif" text-anchor="middle">15</text>
                      </svg>
                    </button>
                    
                    <!-- Play/Pause -->
                    <button 
                      class="bg-red-600 rounded-full p-3 hover:bg-red-700 transition"
                      on:click={togglePlay}
                      aria-label={isPlaying ? 'Pause' : 'Play'}
                    >
                      {#if isPlaying}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z"/>
                        </svg>
                      {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M8 5v14l11-7z"/>
                        </svg>
                      {/if}
                    </button>
                    
                    <!-- Skip forward 30s -->
                    <button 
                      class="text-white hover:text-red-500 transition"
                      on:click={() => {
                        if (audioElement) {
                          audioElement.currentTime = Math.min(duration, audioElement.currentTime + 30);
                          currentTime = audioElement.currentTime;
                          logDebug("Skipped forward 30s");
                        }
                      }}
                      aria-label="Skip forward 30 seconds"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 5V3l4 4-4 4V9c-3.31 0-6 2.69-6 6s2.69 6 6 6 6-2.69 6-6h2c0 4.42-3.58 8-8 8s-8-3.58-8-8 3.58-8 8-8z"/>
                        <text x="12" y="15" font-size="6" font-weight="bold" font-family="sans-serif" text-anchor="middle">30</text>
                      </svg>
                    </button>
                  </div>
                  
                  <div class="flex items-center space-x-4">
                    <!-- Volume control -->
                    <div class="relative">
                      <button 
                        class="text-white hover:text-red-500 transition"
                        on:click={() => showVolumeControl = !showVolumeControl}
                        aria-label="Volume control"
                      >
                        {#if volume === 0}
                          <!-- Muted -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M16.5 12c0-1.77-1.02-3.29-2.5-4.03v2.21l2.45 2.45c.03-.2.05-.41.05-.63zm2.5 0c0 .94-.2 1.82-.54 2.64l1.51 1.51C20.63 14.91 21 13.5 21 12c0-4.28-2.99-7.86-7-8.77v2.06c2.89.86 5 3.54 5 6.71zM4.27 3L3 4.27 7.73 9H3v6h4l5 5v-6.73l4.25 4.25c-.67.52-1.42.93-2.25 1.18v2.06c1.38-.31 2.63-.95 3.69-1.81L19.73 21 21 19.73l-9-9L4.27 3zM12 4L9.91 6.09 12 8.18V4z"/>
                          </svg>
                        {:else if volume < 0.5}
                          <!-- Low volume -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M7 9v6h4l5 5V4l-5 5H7z"/>
                          </svg>
                        {:else}
                          <!-- High volume -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
                          </svg>
                        {/if}
                      </button>
                      
                      {#if showVolumeControl}
                        <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 w-32 p-2 bg-gray-700 rounded shadow-lg">
                          <input 
                            type="range" 
                            min="0" 
                            max="1" 
                            step="0.01" 
                            value={volume} 
                            on:input={handleVolumeChange}
                            class="w-full h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer"
                          />
                        </div>
                      {/if}
                    </div>
                    
                    <!-- Playback speed -->
                    <div class="relative">
                      <button 
                        class="text-white hover:text-red-500 transition px-2 py-1 rounded border border-gray-600"
                        on:click={() => showPlaybackSpeedControl = !showPlaybackSpeedControl}
                        aria-label="Playback speed"
                      >
                        {playbackRate}x
                      </button>
                      
                      {#if showPlaybackSpeedControl}
                        <div class="absolute bottom-full right-0 mb-2 bg-gray-700 rounded shadow-lg overflow-hidden z-10">
                          {#each [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2] as rate}
                            <button 
                              class="block w-full px-4 py-2 text-left hover:bg-gray-600 transition {playbackRate === rate ? 'bg-gray-600' : ''}"
                              on:click={() => setPlaybackRate(rate)}
                            >
                              {rate}x
                            </button>
                          {/each}
                        </div>
                      {/if}
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Reading progress status -->
              <div class="mt-6 text-gray-400 text-sm flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {#if duration > 0}
                  <span>
                    {Math.floor((currentTime / duration) * 100)}% complete
                    {#if isSavingProgress}
                      <span class="ml-2 italic">(Saving progress...)</span>
                    {/if}
                  </span>
                {:else}
                  <span>
                    Duration not available
                    <button 
                      on:click={retryLoadAudio} 
                      class="ml-2 text-red-500 hover:underline"
                    >
                      Retry
                    </button>
                  </span>
                {/if}
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/if}
  </div>