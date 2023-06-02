import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  HStack,
} from '@chakra-ui/react';

interface NewCastModalProps {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
};

const NewCastModal: React.FC<NewCastModalProps> = (props) => {

  const { isOpen, onClose } = props;

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>New Cast</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
        </ModalBody>
        <ModalFooter alignSelf={"center"}>
          <HStack w="full">
            <Button w="full" colorScheme="green">Submit</Button>
            <Button w="full" onClick={onClose}>Close</Button>
          </HStack>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default NewCastModal;
