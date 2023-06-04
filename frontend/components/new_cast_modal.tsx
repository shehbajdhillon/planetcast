import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  Text,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  HStack,
  FormControl,
  FormLabel,
  useColorModeValue,
  Input,
  Stack,
  Box,
  Center,
} from '@chakra-ui/react';
import { Check } from 'lucide-react';
import { useEffect, useState } from 'react';

import Dropzone from "react-dropzone";

interface NewCastModalProps {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
};

const NewCastModal: React.FC<NewCastModalProps> = (props) => {

  const { isOpen, onClose } = props;

  const [castTitle, setCastTitle] = useState("");
  const [mediaFile, setMediaFile] = useState<any>();

  const fileToDataUri = (file: File) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target)
        setMediaFile(event.target.result);
    };
    reader.readAsDataURL(file);
  };


  useEffect(() => {
    setCastTitle("");
    setMediaFile(undefined);
  }, [isOpen, onClose]);

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
              <Dropzone onDrop={(acceptedFiles) => setMediaFile(acceptedFiles[0])}>
                {({ getRootProps, getInputProps }) => (
                  <Box
                    w="full"
                    h="full"
                    minH={"80px"}
                    borderWidth={"2px"}
                    borderRadius={"10px"}
                    borderStyle={"dotted"}
                    {...getRootProps({className: 'dropzone'})}
                  >
                    <Center>
                      <input {...getInputProps()} />
                      <Text>
                        { !mediaFile ? "Upload Media File" : (mediaFile as File).name }
                      </Text>
                    </Center>
                  </Box>
                )}
              </Dropzone>
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