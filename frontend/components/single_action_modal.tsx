import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Center,
} from '@chakra-ui/react'
import NProgress from 'nprogress';
import { useEffect } from 'react';

import Button from './button';


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
              px="16px"
              mr={3}
              flip
              onClick={takeAction}
              isDisabled={loading}
            >
              Submit
            </Button>
            <Button onClick={onClose} px="16px">
              Close
            </Button>
          </Center>

        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default SingleActionModal;
