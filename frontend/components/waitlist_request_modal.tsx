import {
  Box,
  Button,
  useColorModeValue,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Center,
  Input,
  useToast
} from "@chakra-ui/react";
import { useEffect, useState } from "react";

const WaitlistModal: React.FC = () => {

  const { isOpen, onOpen, onClose } = useDisclosure();

  const [email, setEmail] = useState("");

  const [loading, setLoading] = useState(false);

  const toast = useToast();

  function validateEmail(email: string) {
    if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email)) {
      return true;
    }
    return false;
  }

  useEffect(() => {
    setEmail("");
  }, [isOpen, onOpen, onClose]);

  const submitRequest = () => {

    if (loading) return;

    setLoading(true);

    if (!validateEmail(email)) {

      toast({
        title: 'Invalid Email Address.',
        description: "Please check the email entered and try again!",
        status: 'error',
        duration: 6000,
        isClosable: true,
      });

      setLoading(false);

    } else {

      const data: Record<string, any> = {};

      data.waitlist_id = 9227;
      data.referral_link = document.URL;
      data.email = email;

      fetch("https://api.getwaitlist.com/api/v1/waiter", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
      })
      .then((response) => response.json())
      .then((_) => {
        toast({
          title: 'Successfully Submitted Request.',
          description: "You'll be hearing from us soon!",
          status: 'success',
          duration: 6000,
          isClosable: true,
        });
        onClose();
        setLoading(false);
      })
      .catch((_) => {
        toast({
          title: 'Something Went Wrong.',
          description: "Please check the email entered and try again!",
          status: 'error',
          duration: 6000,
          isClosable: true,
        });
        setLoading(false);
      });
    }
  };

  return (
    <Box>
      <Button
        size={"lg"}
        variant={"outline"}
        backgroundColor={useColorModeValue("black", "white")}
        textColor={useColorModeValue("white", "black")}
        onClick={onOpen}
        _hover={{
          backgroundColor: useColorModeValue("black", "white:"),
          textColor: useColorModeValue("white", "black"),
          bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
        }}
      >
        Request Access
      </Button>

      <Modal isOpen={isOpen} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Request Access</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Input
              placeholder="hello@planetcast.ai"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </ModalBody>
          <ModalFooter w="full">
            <Center w={"full"}>
              <Button
                mr={3}
                variant={"outline"}
                backgroundColor={useColorModeValue("black", "white")}
                textColor={useColorModeValue("white", "black")}
                _hover={{
                  backgroundColor: useColorModeValue("black", "white:"),
                  textColor: useColorModeValue("white", "black"),
                  bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
                }}
                onClick={submitRequest}
                isDisabled={loading}
              >
                Submit
              </Button>
              <Button variant={"outline"} onClick={onClose}>
                Close
              </Button>
            </Center>
          </ModalFooter>
        </ModalContent>
      </Modal>

    </Box>
  );

};

export default WaitlistModal;
