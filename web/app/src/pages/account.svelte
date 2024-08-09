<script>
  import { currentUser } from '../lib/userstore.js';
  import { Link } from 'svelte-routing';
  import { Card, Button } from 'flowbite-svelte';
  import { DotsVerticalOutline } from 'flowbite-svelte-icons';
  import { Dropdown, DropdownItem } from 'flowbite-svelte';
  import Transfer from '../assets/Transfer.png';
  import Payment from '../assets/Pay.png';
  import Loan from '../assets/Loan.png';
  import Investment from '../assets/Invest.png';
  import Statement from '../assets/Statement.png';

  let user = null;

  // Subscribe to currentUser store
  currentUser.subscribe(value => {
    user = value;
  });
</script>

<div class="flex flex-col items-center">
  <Card class="bg-white mb-5" size="lg" padding="xl" style="width: 900px;">
    <div class="flex flex-col h-full">
      <div class="flex items-center justify-between mb-4">
        <div class="flex flex-col">
          <h5 class="text-[#004D00] mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            {user ? user.fullname : 'Guest'}
          </h5>
          <h6 class="mb-3 font-normal text-base text-gray-700 dark:text-gray-400 leading-tight">
            {user && user.accountNumber ? user.accountNumber : 'account number'}
          </h6>
        </div>
        <div class="flex items-center">
          <DotsVerticalOutline class="dots-menu dark:text-white cursor-pointer mb-10" />
          <Dropdown triggeredBy=".dots-menu" class="bg-slate-50 rounded shadow-lg">
            <DropdownItem class="bg-white hover:bg-slate-50">Account 1</DropdownItem>
            <DropdownItem class="bg-white hover:bg-slate-50">Account 2</DropdownItem>
            <DropdownItem class="bg-white hover:bg-slate-50">Account 3</DropdownItem>
            <DropdownItem slot="footer" class="bg-white hover:bg-slate-50">
              <Link to="/addaccount" class="w-full text-left block text-gray-700 dark:text-gray-400 hover:bg-slate-50">
                Add new account
              </Link>
            </DropdownItem>
          </Dropdown>
        </div>
      </div>
      <div class="flex flex-1 items-center justify-center">
        <div class="flex flex-col items-center mt-[-80px]">
          <div class="flex flex-col items-center justify-center w-60 h-60 bg-[#28A745] rounded-full border-2 border-gray-300">
            <span class="text-xl font-bold text-white block">Available Balance</span>
            <span class="text-xl font-medium text-white block mt-4">฿1234.56</span>
          </div>
        </div>
      </div>
    </div>
  </Card>
  <div class="flex justify-center space-x-7 mt-8">
    <Link to="/transfer">
      <Button class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2">
        <img src="{Transfer}" alt="Transfer" class="w-8 h-8" />
        <span>Transfer</span>
      </Button>
    </Link>
    <Link to="/payment">
      <Button class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2">
        <img src="{Payment}" alt="Payment" class="w-8 h-8" />
        <span>Payment</span>
      </Button>
    </Link>
    <Link to="/loan">
      <Button class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2">
        <img src="{Loan}" alt="Loan" class="w-8 h-8" />
        <span>Loan</span>
      </Button>
    </Link>
    <Link to="/investment">
      <Button class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2">
        <img src="{Investment}" alt="Investment" class="w-8 h-8" />
        <span>Investment</span>
      </Button>
    </Link>
    <Link to="/statement">
      <Button class="w-40 h-16 bg-white hover:bg-slate-50 text-black flex items-center justify-center space-x-2">
        <img src="{Statement}" alt="Statement" class="w-8 h-8" />
        <span>Statement</span>
      </Button>
    </Link>
  </div>
</div>
