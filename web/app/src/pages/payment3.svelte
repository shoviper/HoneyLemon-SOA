<script>
    import { Link } from "svelte-routing";
    import { Card, Button, Label, Input, Checkbox } from 'flowbite-svelte';
    import HoneyLemonLogo from '../assets/BankLogo.png';
    import CancelButton from '../assets/cancel.png';
    import ConfirmButton from '../assets/confirm.png';
    import NextButton from '../assets/next.png';
    import BackButton from '../assets/back.png';
    import Succ from '../assets/successful.png';
</script>

<Card>
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex flex-col items-center">
        <img src="{Succ}" class="h-18 w-28" alt="HoneyLemonLogo" />
        <span class="text-black text-xl mt-2">Successful</span>
        <span class="text-gray-400 text-sm mt-2">01 August 2024 - 11:34</span>
    </div>
    <div class="flex items-center justify-between">
        <Label class="space-y-2">
            <span class="text-black text-xl">From:</span>
        </Label>
        <Label class="space-y-2 flex flex-col mt-8">
            <span class="text-xl text-[#28A745]">Shogun</span>
            <span class="text-base text-[#666666]">xxx-xxx-xxx</span>
        </Label>
    </div>
    <div class="flex items-center justify-between">
        <Label class="space-y-2">
            <span class="text-black text-xl">To:</span>
        </Label>
        <Label class="space-y-2 flex flex-col mt-8">
            <span class="text-xl text-[#28A745]">AIS</span>
            <span class="text-base text-[#666666]">xxx-xxx-xxx</span>
        </Label>
    </div>
    <div class="flex items-center justify-between">
        <Label class="space-y-2">
            <span class="text-black text-xl">Amount:</span>
        </Label>
        <Label class="space-y-2 ">
            <span class="text-black text-xl ">5000.00</span>
        </Label>
    </div>
    <div class="flex items-center justify-between">
        <Label class="space-y-2">
            <span class="text-black text-xl">Balance:</span>
        </Label>
        <Label class="space-y-2 ">
            <span class="text-black text-xl ">10.00</span>
        </Label>
    </div>
    <Link to="/">
        <Button type="submit" class="w-full bg-[#218838] hover:bg-[#28A745]">
            Back to Account page
        </Button>
    </Link>
  </form>
</Card>
