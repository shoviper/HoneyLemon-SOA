<script>
    import { Link, navigate } from "svelte-routing";
    import { Card, Button, Label, Input, Checkbox, Modal } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons';
    import HoneyLemonLogo from '../assets/BankLogo.png';
    import CancelButton from '../assets/cancel.png';
    import ConfirmButton from '../assets/confirm.png';
    import NextButton from '../assets/next.png';
    import BackButton from '../assets/back.png';

    let popupModal_cancel = false;

    function handleCancelClick() {
        popupModal_cancel = true;
    }

    function handleTransactionCancel() {
        navigate('/mainaccount');
    }
</script>

<Card>
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex items-center justify-between">
        <Link to="/mainaccount">
            <img src="{BackButton}" class="h-4 w-4" alt="BackButton" />
        </Link>
        <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-black text-xl">To:</span>
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400 ml-4">Contact account / Mobile no.</span>
      <Input type="text" name="accno" required class="ml-4 w-11/12"/>
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400 ml-4">Amount</span>
      <Input type="text" name="amount" required class="ml-4 w-11/12"/>
    </Label>
    <div class="flex items-center justify-between">
        <div class="flex items-center justify-between">
            <img src="{CancelButton}" class="h-12 w-12 mt-7 cursor-pointer" alt="Cancel" on:click={handleCancelClick} />
            <span class="text-black ml-1 mt-7">Cancel</span>
        </div>
        <div class="flex items-center justify-between">
            <span class="text-gray-400 mr-1 mt-7">Next</span>
            <Link to="/payment2">
                <img src="{NextButton}" class="h-12 w-12 mt-7" alt="NextButton" />
            </Link>
        </div>
        
    </div>
  </form>
</Card>
