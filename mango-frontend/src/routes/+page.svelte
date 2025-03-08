<script lang="ts">
  import { goto } from '$app/navigation';
  import "../app.css"
  
  let isLogin = true; // Toggle between login and signup
  let message = "";
  let username = "";
  let password = "";
  let email = "";
  let imageFile: File | null = null;
  let imagePreviewUrl: string | null = null;
 
  // Handle image selection with preview
  function handleFileUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      // Only allow a single file
      const file = target.files[0];
      
      // Check file size (10MB max)
      if (file.size > 10 * 1024 * 1024) {
        message = "File is too large. Maximum size is 10MB.";
        target.value = ""; // Reset the input
        return;
      }
      
      // Create a preview URL - this is where the error was happening
      imageFile = file;
      imagePreviewUrl = URL.createObjectURL(file);
    }
  }
  
  // Clean up object URLs when component is destroyed
  import { onDestroy } from 'svelte';
  onDestroy(() => {
    if (imagePreviewUrl) {
      URL.revokeObjectURL(imagePreviewUrl);
    }
  });
  
  // Function to remove the image
  function removeImage() {
    imageFile = null;
    if (imagePreviewUrl) {
      URL.revokeObjectURL(imagePreviewUrl);
      imagePreviewUrl = null;
    }
    // Reset the file input
    const fileInput = document.getElementById('file-upload') as HTMLInputElement;
    if (fileInput) {
      fileInput.value = "";
    }
  }
 
  async function submitForm() {
    if (!username || !password || (!isLogin && !imageFile)) {
      message = "Please fill in all fields.";
      return;
    }
   
    const formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);
    
    if (!isLogin && imageFile){
      formData.append("image", imageFile); // Only for signup
      formData.append("email", email)
    } 
   
    const endpoint = isLogin ? "login" : "signup";
    
    try {
      const res = await fetch(`https://api.thejoyestboy.com/${endpoint}`, {
        method: "POST",
        body: formData,
        credentials: "include" // Store session cookies
      });
     
      if (res.ok) {
        if (isLogin) {
          // If login is successful, navigate to home
          message = "Login successful!";
          setTimeout(() => goto('/home'), 50);
        } else {
          // If signup is successful, show message and switch to login view
          message = "Signup successful! Please log in with your new account.";
          
          // Clear form fields
          username = "";
          password = "";
          email = "";
          
          // Clear image
          removeImage();
          
          // Switch to login view after 1.5 seconds
          setTimeout(() => {
            isLogin = true;
            message = ""; // Clear the message when switching to login view
          }, 1500);
        }
      } else {
        const errorData = await res.json().catch(() => null);
        message = errorData?.message || `Request failed with status ${res.status}`;
      }
    } catch (error) {
      console.error("Error during form submission:", error);
      message = "Network error. Please try again.";
    }
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-gradient-to-br from-indigo-100 to-purple-100">
  <div class="w-full max-w-md p-8 space-y-8 bg-white rounded-xl shadow-lg">
    <!-- Card Header -->
    <div class="text-center">
      <h2 class="text-3xl font-extrabold text-gray-900">
        {isLogin ? "Welcome Back" : "Create Account"}
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        {isLogin ? "Sign in to access your account" : "Sign up to get started"}
      </p>
    </div>
    
    <!-- Form Area -->
    <div class="mt-8 space-y-6">
      <!-- Error/Success Message -->
      {#if message}
        <div class={`p-4 rounded-lg text-sm font-medium text-center ${message.includes('successful') ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'}`}>
          {message}
        </div>
      {/if}
      
      <!-- Username Field -->
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
        <div class="mt-1">
          <input
            id="username"
            type="text"
            bind:value={username}
            class="block w-full px-3 py-2 placeholder-gray-400 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="Enter your username"
          />
        </div>
      </div>
      
      <!-- Password Field -->
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
        <div class="mt-1">
          <input
            id="password"
            type="password"
            bind:value={password}
            class="block w-full px-3 py-2 placeholder-gray-400 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="Enter your password"
          />
        </div>
      </div>
      
      
      <!-- Image Upload (Signup only) -->
      {#if !isLogin}
          <!-- Email Field -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
            <div class="mt-1">
              <input
                id="email"
                type="email"
                bind:value={email}
                class="block w-full px-3 py-2 placeholder-gray-400 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Enter your email"
              />
            </div>
          </div>
          <div>
            <label for="image" class="block text-sm font-medium text-gray-700">Invite Image</label>
            {#if imagePreviewUrl}
              <!-- Display preview of the uploaded image -->
              <div class="mt-1 flex flex-col items-center justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md">
                <div class="mx-auto h-32 w-32 overflow-hidden rounded-full border-2 border-indigo-500">
                  <img src={imagePreviewUrl} alt="Profile preview" class="h-full w-full object-cover" />
                </div>
                <button
                  type="button"
                  class="mt-4 text-sm text-red-600 hover:text-red-700"
                  on:click={removeImage}
                >
                  Remove image
                </button>
              </div>
            {:else}
              <!-- Entire box is clickable -->
              <label for="file-upload" class="mt-1 flex flex-col items-center justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md cursor-pointer hover:border-indigo-500 hover:bg-indigo-50 transition-colors">
                <div class="space-y-1 text-center">
                  <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48" aria-hidden="true">
                    <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                  <div class="text-sm text-gray-600">
                    <span class="font-medium text-indigo-600 hover:text-indigo-500">Upload a file</span>
                    <span class="pl-1">or drag and drop</span>
                  </div>
                  <p class="text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
                </div>
                <input id="file-upload" name="file-upload" type="file" accept="image/*" on:change={handleFileUpload} class="sr-only" />
              </label>
            {/if}
          </div>
      {/if}
      
      
      <!-- Submit Button -->
      <div>
        <button
          on:click={submitForm}
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition duration-150"
        >
          {isLogin ? "Sign In" : "Create Account"}
        </button>
      </div>
      
        <!-- Toggle Login/Signup -->
        <div class="text-sm text-center">
          <button 
            type="button"
            on:click={() => isLogin = !isLogin} 
            class="font-medium text-indigo-600 hover:text-indigo-500 cursor-pointer border-none bg-transparent p-0"
          >
            {isLogin ? "Don't have an account? Sign up" : "Already have an account? Sign in"}
          </button>
      </div>
    </div>
  </div>
</div>