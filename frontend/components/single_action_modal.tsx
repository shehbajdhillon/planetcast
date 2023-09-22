import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  Center,
  useColorModeValue,
} from '@chakra-ui/react'
import NProgress from 'nprogress';
import { useEffect } from 'react';


interface SingleActionModalProps {
  heading: any;
  body: any;
  action: () => Promise<void>;
  isOpen: boolean;
  onClose: () => void;
  loading?: boolean;
};


const SingleActionModal: React.FC<SingleActionModalProps> = (props) => {

  const {heading, body, action, isOpen, onClose, loading} = props;

  const takeAction = async () => {
    await action();
    onClose();
  };

  useEffect(() => {
    if (loading) NProgress.start();
    else NProgress.done();
  }, [loading]);

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{heading}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          {body}
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
              onClick={takeAction}
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
  );
};

export default SingleActionModal;
