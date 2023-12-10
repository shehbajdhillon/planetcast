import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Center,
  ModalProps,
} from '@chakra-ui/react'
import NProgress from 'nprogress';
import { ReactNode, useEffect } from 'react';

import Button from './button';


interface SingleActionModalProps extends ModalProps {
  heading: any;
  action: () => any;
  isOpen: boolean;
  onClose: () => void;
  loading?: boolean;
  disabled?: boolean;
  children: ReactNode;
};


const SingleActionModal: React.FC<SingleActionModalProps> = (props) => {

  const { disabled, heading, children, action, onClose, loading } = props;

  const takeAction = async () => {
    await action();
    onClose();
  };

  useEffect(() => {
    if (loading) NProgress.start();
    else NProgress.done();
  }, [loading]);

  return (
    <Modal isCentered {...props}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{heading}</ModalHeader>
        <ModalCloseButton />
        <ModalBody overflowY={"auto"} maxH={"600px"}>
          {children}
        </ModalBody>
        <ModalFooter w="full">
          <Center w={"full"}>
            <Button
              px="16px"
              mr={3}
              flip
              onClick={takeAction}
              isDisabled={loading || disabled}
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
