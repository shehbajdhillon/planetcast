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
  FormControl,
  FormLabel,
  Icon,
  useColorModeValue,
  Input,
  Stack,
  Spacer,
  Center,
} from '@chakra-ui/react';
import { BadgeCheck, Check, Info } from 'lucide-react';
import { useEffect, useState } from 'react';

interface NewCastModalProps {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
};

const NewCastModal: React.FC<NewCastModalProps> = (props) => {

  const { isOpen, onClose } = props;

  const [castTitle, setCastTitle] = useState("");
  const [mediaFile, setMediaFile] = useState<any>();

  useEffect(() => {
    console.log({ mediaFile });
  }, [mediaFile]);

  const fileToDataUri = (file: File) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target)
        setMediaFile(event.target.result);
    };
    reader.readAsDataURL(file);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered size={{ base: "lg", md: "2xl" }}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>New Cast</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FormControl isRequired={true}>
            <Stack spacing={-1} py="10px">
              <FormLabel
                textColor={useColorModeValue('gray.700', 'white')}
                fontWeight={'600'}
                fontSize={'lg'}
              >
                Title
              </FormLabel>
              <Input
                placeholder="PlanetCast Episode #234"
                onChange={(e) => setCastTitle(e.target.value)}
                value={castTitle}
              />
            </Stack>
            <Stack spacing={-1} py="10px">
              <FormLabel
                textColor={useColorModeValue('gray.700', 'white')}
                fontWeight={'600'}
                fontSize={'lg'}
              >
                Media File
              </FormLabel>
              <Input
                type='file'
                colorScheme='teal'
                borderWidth={"0px"}
                onChange={(e) => {
                  if (e.target.files)
                    fileToDataUri(e.target.files[0])
                }}
              />
            </Stack>
          </FormControl>
        </ModalBody>
        <ModalFooter alignSelf={"center"}>
          <HStack w="full">
            <Button
              leftIcon={<Check />}
              w="full"
              colorScheme="green"
              onClick={onClose}
              px="25px"
            >
              Submit
            </Button>
            <Button
              w="full"
              onClick={onClose}
              px="25px"
            >
              Close
            </Button>
          </HStack>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default NewCastModal;
