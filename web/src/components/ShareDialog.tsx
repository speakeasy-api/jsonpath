import {
  Button,
  Dialog,
  Heading,
  Separator,
  Stack,
  Text,
} from "@speakeasy-api/moonshine";
import { forwardRef, useImperativeHandle, useState } from "react";
import { CopyButton } from "./CopyButton";

export interface ShareDialogHandle {
  setUrl: React.Dispatch<React.SetStateAction<string>>;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const ShareDialog = forwardRef<ShareDialogHandle, {}>((_, ref) => {
  const [url, setUrl] = useState<string>("");
  const [open, setOpen] = useState(false);

  useImperativeHandle(ref, () => ({
    setUrl,
    setOpen,
  }));

  const handleClose = () => {
    setOpen(false);
  };

  const handleOpenChange = (open: boolean) => {
    setOpen(open);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <Dialog.Content>
        <Dialog.Header>
          <Dialog.Title asChild>
            <div>
              <Heading variant="lg">Share</Heading>
              <Text muted variant="sm" className="leading-none ">
                Copy and paste the URL below anywhere to share this overlay
                session with others.
              </Text>
            </div>
          </Dialog.Title>
        </Dialog.Header>
        <Separator />
        <Stack direction="vertical" gap={10} className="my-2">
          <CopyButton className="w-full" value={url} />
        </Stack>
        <Separator />
        <Dialog.Footer>
          <Dialog.Close asChild>
            <Button onClick={handleClose}>Done</Button>
          </Dialog.Close>
        </Dialog.Footer>
      </Dialog.Content>
    </Dialog>
  );
});
ShareDialog.displayName = "ShareDialog";

export default ShareDialog;
