<script>
    import { Link, navigate } from "svelte-routing";
    import { Card, Button, Label, Input } from 'flowbite-svelte';
    import { ArrowRightOutline } from 'flowbite-svelte-icons';
    import { Dropdown, DropdownItem } from 'flowbite-svelte';
    import { DotsVerticalOutline } from 'flowbite-svelte-icons';
    import PaymentPic from '../assets/Payment.png';
    import Payment from '../assets/Pay.png';
    import Loan from '../assets/Loan.png';
    import Transfer from '../assets/Transfer.png';
    import Statement from '../assets/Statement.png';
    import Investment from '../assets/Invest.png';
    import HoneyLemonLogo from '../assets/BankLogo.png';
    import BackButton from '../assets/back.png';
</script>

<div class="flex flex-col items-center">
  <Card class="bg-white mb-5" size="md" padding="xl" style="width: 900px;">
    <div class="flex flex-col h-full">
        <div class="flex items-center justify-between">
            <Link to="/">
                <img src="{BackButton}" class="h-4 w-4" alt="HoneyLemonLogo" />
            </Link>
            <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
        </div>
        <div class="text-base text-black mt-4 ml-4 mb-3">
            Request Statement
        </div>
        <div class="ml-4">
            <Card class="bg-white border-black mb-5" size="md" padding="xl">
                <div class="flex flex-col h-full">
                    <div class="text-base text-black">
                        Select channel
                    </div>
                    <div class="flex items-center space-x-4 mt-2">
                        <Label class="inline-flex items-center">
                        <input type="radio" name="channelway" value="HoneyLemon"  class="appearance-none h-4 w-4 border border-gray-300 rounded-full checked:bg-green-600 checked:border-transparent focus:outline-none transition duration-200" />
                        <span class="ml-2 text-gray-700 dark:text-gray-300">HoneyLemon</span>
                        </Label>
                        <Label class="inline-flex items-center">
                        <input type="radio" name="channelway" value="E-mail"  class="appearance-none h-4 w-4 border border-gray-300 rounded-full checked:bg-green-600 checked:border-transparent focus:outline-none transition duration-200" />
                        <span class="ml-2 text-gray-700 dark:text-gray-300">E-mail</span>
                        </Label>
                    </div>
                </div>
                <div class="flex flex-col h-full">
                    <div class="text-base text-black mt-2">
                        Start date
                    </div>
                    <Input type="date" name="startdate" placeholder="yyyy-mm-dd" required class="w-2/3 mt-2"/>
                </div>
                <div class="flex flex-col h-full">
                    <div class="text-base text-black mt-2">
                        End date
                    </div>
                    <Input type="date" name="enddate" placeholder="yyyy-mm-dd" required class="w-2/3 mt-2"/>
                </div>  
            </Card>
        </div>
        <Button type="submit" class="w-full bg-[#218838] hover:bg-[#28A745]">
            Submit request
        </Button>
    </div>
  </Card>
</div>