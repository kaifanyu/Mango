<script>
    import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Button, Input } from 'flowbite-svelte';
    import { SearchOutline } from 'flowbite-svelte-icons';
    import "../app.css";
  </script>

<Navbar class="bg-gray-800 text-white">
  <NavBrand href="/">
      <img src="/favicon.png" class="me-3 h-6 sm:h-9" alt="Flowbite Logo" />
      <span class="self-center whitespace-nowrap text-xl font-semibold">WEBSITE</span>
  </NavBrand>

  <div class="flex md:order-2 ml-80">
      <Button color="none" data-collapse-toggle="mobile-menu-3" aria-controls="mobile-menu-3" aria-expanded="false"
          class="md:hidden text-gray-200 hover:bg-gray-700 focus:outline-hidden focus:ring-4 focus:ring-gray-600 rounded-lg text-sm p-2.5 me-1">
          <SearchOutline class="w-5 h-5" />
      </Button>
      <div class="hidden relative md:block">
          <div class="flex absolute inset-y-0 start-0 items-center ps-3 pointer-events-none">
              <SearchOutline class="w-4 h-4 text-gray-400" />
          </div>
          <Input id="search-navbar" class="ps-10 text-black" placeholder="Search..." />
      </div>
      <NavHamburger />
  </div>

  <NavUl class="md:flex md:space-x-8 bg-gray-800">
      <NavLi 
      href="/" 
      activeClass="text-white font-semibold border-b-2 border-white"
      nonActiveClass="text-gray-400 hover:text-gray-300"
    >
      Home
    </NavLi>

      <NavLi href="/about" class="text-white hover:text-gray-300">About</NavLi>
      <NavLi href="/docs/components/navbar" class="text-white hover:text-gray-300">Navbar</NavLi>
  </NavUl>
</Navbar>


<!-- Page Content -->
<main>
  <slot />
</main>