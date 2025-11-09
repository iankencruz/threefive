<script lang="ts">
	let { variant = "opaque" } = $props();

	let NavLinks = [
		{ id: 1, title: "Home", href: "/" },
		{ id: 2, title: "About Us", href: "/about" },
		{ id: 3, title: "Contact Us", href: "/contact" },
	];
</script>

<nav class="py-10 w-full absolute  transition-all z-20 duration-500">
    <div class="mx-auto max-w-8xl px-4 sm:px-6 lg:px-8">
      <div class="w-full flex flex-col lg:flex-row">
        <div class="flex justify-between lg:flex-row ">
          <a href="/" class="flex items-center w-max">
          Threefive Project
          </a>
          <div class="flex items-center justify-end gap-5">
            <button data-collapse-toggle="navbar" type="button"
              class="inline-flex items-center p-2 text-sm text-white hover:text-gray-900 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200"
              aria-controls="navbar-default" aria-expanded="false">
              <span class="sr-only">Open main menu</span>
              <svg class="w-6 h-6" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20"
                xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd"
                  d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                  clip-rule="evenodd"></path>
              </svg>
            </button>
          </div>

        </div>
        <div class="hidden w-full lg:flex lg:pl-11 max-lg:mt-1 max-lg:overflow-y-auto" id="navbar">
          <ul
            class="flex lg:items-center lg:justify-center flex-col gap-6 max-lg:pt-4 max-lg:mb-4 lg:mt-0 lg:flex-row lg:ml-auto ">
          {#each NavLinks as link (link.id)}
            <li>
              <a href={link.href}
                class="nav-link block lg:mx-3 md:mb-0 lg:text-left text-base text-white font-medium transition-all duration-500 hover:text-gray-400">
                {link.title}
              </a>
            </li>
          {/each}
            
            <li class="relative">
              <button data-target="menu2"
                class="dropdown-toggle flex items-center justify-between text-white text-center text-base font-medium transition-all duration-500 lg:mx-3 lg:mb-0 mr-auto lg:text-left lg:m-0 hover:text-gray-400">
                Features
                <svg class="w-3 h-2 ml-1.5" width="10" height="6" viewBox="0 0 10 6" fill="none"
                  xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M1 1L3.58579 3.58579C4.25245 4.25245 4.58579 4.58579 5 4.58579C5.41421 4.58579 5.74755 4.25245 6.41421 3.58579L9 1"
                    stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"></path>
                </svg>
              </button>
              <!-- Dropdown menu -->
              <div id="menu2" aria-labelledby="menu2"
                class="dropdown-menu z-10 relative top-3 max-lg:mb-3 lg:absolute lg:top-14 font-normal bg-white rounded-lg w-64 xl:p-8 lg:p-4 p-2 hidden">
                <ul class="text-sm text-gray-700" aria-labelledby="dropdownLargeButton">
                  <h6 class="text-sm text-gray-500 font-medium mb-2">
                    Features
                  </h6>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Notification</a>
                  </li>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Analytics</a>
                  </li>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Integrations</a>
                  </li>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Security</a>
                  </li>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Documentation</a>
                  </li>
                  <li>
                    <a href="javascript:;"
                      class="block py-3 hover:text-prime-blue-700 text-base text-gray-900 font-semibold transition-all duration-500">Support</a>
                  </li>
                </ul>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
</nav>
