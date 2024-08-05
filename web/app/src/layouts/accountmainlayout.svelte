<script>
    import { Link, navigate } from "svelte-routing";
    import { Modal } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons';
    import HoneyLemonLogo from "../assets/BankLogo.png";
    import AccountLogo from "../assets/AccountLogo.png";
    import { Button, Dropdown, DropdownItem, Avatar, DropdownHeader, DropdownDivider } from 'flowbite-svelte';
    import { BellSolid, EyeSolid } from 'flowbite-svelte-icons';

    let popupModal = false;

    function handleConfirmClick() {
        navigate('/index');
    }
</script>

<div class="navbar bg-[#F0F0F0] h-24 flex items-center justify-between px-5 fixed top-0 left-0 right-0 shadow-lg z-50">
    <div class="flex flex-row items-center">
        <div class="ml-1">
        <img src="{HoneyLemonLogo}" class="Logo" alt="HoneyLemonLogo" />
        </div>
    </div>
    <div>
        <Avatar class="acs mr-12" src="{AccountLogo}" />
        <Dropdown class="bg-[#F0F0F0] rounded shadow-lg" triggeredBy=".acs">
            <DropdownHeader class="bg-[#28A745] p-4 rounded-t-lg">
                <span class="block text-sm text-white dark:text-white">Inthat Sappipat</span>
            </DropdownHeader>
            <DropdownDivider class="my-2 border-t border-gray-300" />
            <DropdownItem class="bg-white hover:bg-slate-50 px-4 py-2 text-gray-700">
              <Link to="/" class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50">
                  Home
              </Link>
            </DropdownItem>
            <DropdownItem class="bg-white hover:bg-slate-50 px-4 py-2 text-gray-700">
              <Link to="/changeprofile" class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50">
                  Change profile
              </Link>
            </DropdownItem>
            <DropdownItem class="bg-white hover:bg-slate-50 px-4 py-2 text-gray-700">
              <Link to="/changepassword" class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50">
                  Change password
              </Link>
            </DropdownItem>
            <DropdownItem class="bg-white hover:bg-slate-50 px-4 py-2 text-gray-700">
              <Link to="/changepin" class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50">
                  Change pin
              </Link>
            </DropdownItem>
            <DropdownDivider class="my-2 border-t border-gray-300" />
            <DropdownItem class="bg-white hover:bg-slate-50 px-4 py-2 text-gray-700">
              <div class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50" on:click={() => (popupModal = true)}>
                  Sign out
              </div>
            </DropdownItem>
            <Modal bind:open={popupModal} size="xs" autoclose>
            <div class="text-center">
                <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12" />
                <h3 class="mb-5 text-lg font-normal text-gray-500">Are you sure you want to sign out?</h3>
                <div class="flex justify-center gap-2">
                    <Button color="red" on:click={handleConfirmClick}>Yes, I'm sure</Button>
                    <Button color="alternative" on:click={() => (popupModal = false)}>No, cancel</Button>
                </div>
            </div>
        </Modal>
        </Dropdown>
  </div>
</div>

<div class="content-container pt-24">
  <slot></slot>
</div>

<style>
  .Logo {
    height: 36px;
  }
  .LoginButton {
    background-color: #ffffff;
    color: #000000;
    font-size: 16px;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  .LoginButton:hover {
    opacity: 70%;
  }
  .content-container {
    padding-top: 96px; 
  }
</style>
