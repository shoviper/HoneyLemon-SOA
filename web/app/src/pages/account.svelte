<script>
  import { Card, Button, Dropdown, DropdownItem } from "flowbite-svelte";
  import { Link, navigate } from "svelte-routing";
  import { DotsVerticalOutline } from "flowbite-svelte-icons";
  import { currentUser, users } from "../lib/userstore.js";
  import { onMount } from "svelte";
  import HoneyLemonLogo from "../assets/BankLogo.png";
  import Transfer from "../assets/Transfer.png";
  import Payment from "../assets/Pay.png";
  import Loan from "../assets/Loan.png";
  import Investment from "../assets/Invest.png";
  import Statement from "../assets/Statement.png";
  import axios from "axios";

  let user = null;
  let localUsers = [];
  let userData = null;
  let accountData = null;
  let selectedAccount = null;
  let loggedIn = false; // Default state for logged in status

  function checkLoginStatus() {
    // Check for the presence of a specific cookie
    const cookies = document.cookie.split(";").map((cookie) => cookie.trim());
    const authCookie = cookies.find((cookie) =>
      cookie.startsWith("esb_token=")
    );

    loggedIn = !!authCookie;
  }

  onMount(() => {
    checkLoginStatus();
    if (!loggedIn) {
      navigate("/");
    }

    users.subscribe((value) => {
      localUsers = value;
    });

    currentUser.subscribe((value) => {
      user = value;
    });

    const token = getCookie('esb_token');
    if (token) {
      fetchData(token);
    }
  });

  function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
  }

  async function fetchData(token) {
    try {
      // Set the token as a cookie
      document.cookie = `esb_token=${token}; path=/;`;

      const accResponse = await axios.get('http://127.0.0.1:4000/esb/accounts/clientAcc', {
        withCredentials: true, // Ensure credentials are sent with the request
        headers: {
          'esb_token': `Bearer ${token}`
        }
      });

      const userResponse = await axios.get('http://127.0.0.1:4000/esb/clients/info', {
        withCredentials: true, // Ensure credentials are sent with the request
        headers: {
          'esb_token': `Bearer ${token}`
        }
      });

      accountData = accResponse.data;
      if (accountData == null) {
        console.log('accountData:', accountData);
      }

      userData = userResponse.data;
      // console.log('userData:', userData);

      selectedAccount = accountData.length > 0 ? accountData[0] : null;
      // console.log('selectedAccount:', selectedAccount);
      // Update user accounts with the fetched data
      // user = {
      //   thisuser: userData,
      //   accounts: accountData,
      //   selectedAccount: accountData.length > 0 ? accountData[0].id : null
      // };


      // currentUser.set(user);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }

  function switchAccount(accountNumber) {
    selectedAccount = accountData.find((account) => account.id === accountNumber);
  }
</script>
<div class="flex flex-col items-center">
  <Card class="bg-white mb-5" size="lg" padding="xl" style="width: 900px;">
    <div class="flex flex-col h-full">
      <div class="flex items-center justify-between mb-4">
        <div class="flex flex-col">
          <h5
            class="text-[#004D00] mb-2 text-xl font-bold tracking-tight text-gray-900 dark:text-white"
          >
            {#if userData}
              {userData.name}
            {/if}
          </h5>
          <h6
            class="mb-3 font-normal text-base text-gray-700 dark:text-gray-400 leading-tight"
          >
            {userData && selectedAccount
              ? `${selectedAccount.id || "xxxxx"}`
              : "xxxxx"}
          </h6>
        </div>
        <div class="flex items-center">
          <DotsVerticalOutline
            class="dots-menu dark:text-white cursor-pointer mb-10"
          />
          <Dropdown triggeredBy=".dots-menu" class="bg-slate-100 rounded shadow-lg">
            {#if accountData && accountData.length > 0}
              {#each accountData as account}
                <DropdownItem
                  class="bg-white hover:bg-slate-50 text-gray-700"
                  on:click={() => switchAccount(account.id)}
                >
                  {account.id}
                </DropdownItem>
              {/each}
            {:else}
              <!-- No accounts, show only "Add new account" -->
              <DropdownItem class="bg-white text-gray-500">
                No accounts available
              </DropdownItem>
            {/if}
            <DropdownItem slot="footer" class="bg-[#28A745] hover:bg-[#03C04A]">
              <Link
                to="/addaccount"
                class="w-full text-left block text-white dark:text-gray-400"
              >
                Add new account
              </Link>
            </DropdownItem>
          </Dropdown>
        </div>
      </div>
      <div class="flex flex-1 items-center justify-center">
        <div class="flex flex-col items-center mt-[-80px]">
          <div
            class="flex flex-col items-center justify-center w-60 h-60 bg-[#28A745] rounded-full border-2 border-gray-300"
          >
            <span class="text-xl font-bold text-white block"
              >Available Balance</span
            >
            <span class="text-xl font-medium text-white block mt-4">
              {userData && selectedAccount
                ? `${selectedAccount.balance || 0} ฿`
                : "0 ฿"}
            </span>
          </div>
        </div>
      </div>
    </div>
  </Card>
  <div class="flex justify-center space-x-7 mt-8">
    <Link to="/transfer">
      <Button
        class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2"
      >
        <img src={Transfer} alt="Transfer" class="w-8 h-8" />
        <span>Transfer</span>
      </Button>
    </Link>
    <Link to="/payment">
      <Button
        class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2"
      >
        <img src={Payment} alt="Payment" class="w-8 h-8" />
        <span>Payment</span>
      </Button>
    </Link>
    <Link to="/loan">
      <Button
        class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2"
      >
        <img src={Loan} alt="Loan" class="w-8 h-8" />
        <span>Loan</span>
      </Button>
    </Link>
    <Link to="/investment">
      <Button
        class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2"
      >
        <img src={Investment} alt="Investment" class="w-8 h-8" />
        <span>Investment</span>
      </Button>
    </Link>
    <Link to="/statement">
      <Button
        class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2"
      >
        <img src={Statement} alt="Statement" class="w-8 h-8" />
        <span>Statement</span>
      </Button>
    </Link>
  </div>
</div>
