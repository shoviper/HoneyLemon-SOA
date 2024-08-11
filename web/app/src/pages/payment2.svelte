<script>
    import { Link, navigate } from "svelte-routing";
    import { Card, Input, Label, Button, Modal } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons';
    import HoneyLemonLogo from '../assets/BankLogo.png';
    import CancelButton from '../assets/cancel.png';
    import ConfirmButton from '../assets/confirm.png';
    import BackButton from '../assets/back.png';

    let popupModal_confirm = false;
    let popupModal_cancel = false;

    function handleConfirmClick() {
        popupModal_confirm = true;
    }

    function handleCancelClick() {
        popupModal_cancel = true;
    }

    function handleTransactionConfirm() {
        navigate('/transfer3');
    }

    function handleTransactionCancel() {
        navigate('/mainaccount');
    }
</script>

<Card>
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex items-center justify-between">
        <Link to="/transfer">
            <img src="{BackButton}" class="h-4 w-4" alt="Back" />
        </Link>
        <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
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
        <Label class="space-y-2">
            <span class="text-black text-xl">5000.00</span>
        </Label>
    </div>
    <div class="flex items-center justify-between">
        <Label class="space-y-2 flex flex-col">
            <span class="text-red-500">Beware of scammers!</span>
            <span class="text-gray-400">Please check the recipient's name and account number every time before pressing the confirm button.</span>
        </Label>
    </div>
    <div class="flex items-center justify-between">
        <div class="flex items-center">
            <img src="{CancelButton}" class="h-12 w-12 mt-7 cursor-pointer" alt="Cancel" on:click={handleCancelClick} />
            <span class="text-black ml-1 mt-7">Cancel</span>
        </div>
        <div class="flex items-center">
            <span class="text-black mr-1 mt-7">Confirm</span>
            <img src="{ConfirmButton}" class="h-12 w-12 mt-7 cursor-pointer" alt="Confirm" on:click={handleConfirmClick}/>
        </div>
        <Modal bind:open={popupModal_cancel} size="xs" autoclose>
            <div class="text-center">
                <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12" />
                <h3 class="mb-5 text-lg font-normal text-gray-500">Are you sure you want to cancel?</h3>
                <div class="flex justify-center gap-2">
                    <Button color="red" on:click={handleTransactionCancel}>Yes, I'm sure</Button>
                    <Button color="alternative" on:click={() => (popupModal_cancel = false)}>No, cancel</Button>
                </div>
            </div>
        </Modal>
        <Modal bind:open={popupModal_confirm} size="xs" autoclose noCloseButton>
            <form class="flex flex-col space-y-6" action="#">
                <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Confirm for a pin</h3>
                <Label class="space-y-2">
                  <span>Enter a pin</span>
                  <Input type="text" name="confirmpin" required />
                </Label>
                <div class="flex justify-center gap-64">
                    <Button color="red" on:click={handleTransactionConfirm}>Enter</Button>
                    <Button color="alternative" on:click={() => (popupModal_confirm = false)}>Close</Button>
                </div>
            </form>
        </Modal>
    </div>
  </form>
</Card>
