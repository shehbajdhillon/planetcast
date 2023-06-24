import { SupportedLanguage, SupportedLanguages } from '@/types';
import { gql, useMutation } from '@apollo/client';
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
  Input,
  Stack,
  Box,
  Center,
  Spacer,
  Select,
} from '@chakra-ui/react';
import { Check } from 'lucide-react';
import { useEffect, useState } from 'react';

import Dropzone from "react-dropzone";

const CREATE_PROJECT = gql`
  mutation CreateProject($teamSlug: String!, $title: String!, $sourceLanguage: SupportedLanguage!, $targetLanguage: SupportedLanguage!, $sourceMedia: Upload!) {
    createProject(teamSlug: $teamSlug, title: $title, sourceLanguage: $sourceLanguage, targetLanguage: $targetLanguage, sourceMedia: $sourceMedia) {
      id
      title
    }
  }
`;

interface NewCastModalProps {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
  refetch: () => void;
  teamSlug: string;
};

const NewProjectModal: React.FC<NewCastModalProps> = (props) => {

  const { isOpen, onClose, teamSlug, refetch } = props;

  const [title, setTitle] = useState("");
  const [sourceMedia, setSourceMedia] = useState<File>();
  const [sourceLanguage, setSourceLanguage] = useState<SupportedLanguage>("ENGLISH");
  const [targetLanguage, setTargetLanguage] = useState<SupportedLanguage>("HINDI");

  const [formValid, setFormValid] = useState(false);

  const [createProjectMutation] = useMutation(CREATE_PROJECT);

  const createProject = async () => {
    const res = await createProjectMutation({ variables: { title, teamSlug, sourceLanguage, targetLanguage, sourceMedia } });
    if (res) {
      refetch();
      onClose();
    }
  };

  useEffect(() => {
    const checkFormValid = () => {
      if (!title.length) return false;
      if (!sourceMedia) return false;
      if (sourceLanguage === targetLanguage) return false;
      if (SupportedLanguages.indexOf(sourceLanguage) === -1) return false;
      if (SupportedLanguages.indexOf(targetLanguage) === -1) return false;
      return true;
    };
    setFormValid(checkFormValid());
  }, [title, sourceMedia, sourceLanguage, targetLanguage]);

  useEffect(() => {
    setTitle("");
    setSourceMedia(undefined);
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
                fontWeight={'600'}
                fontSize={'lg'}
              >
                Title
              </FormLabel>
              <Input
                placeholder="PlanetCast Episode #234"
                onChange={(e) => setTitle(e.target.value)}
                value={title}
              />
            </Stack>
            <Stack spacing={-1} py="10px">
              <FormLabel
                fontWeight={'600'}
                fontSize={'lg'}
              >
                Media File
              </FormLabel>
              <Dropzone onDrop={(acceptedFiles) => setSourceMedia(acceptedFiles[0])}>
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
                      <input {...getInputProps()} type='file' />
                      <Text py="20px">
                        { !sourceMedia? "Upload Media File" : (sourceMedia as File).name }
                      </Text>
                    </Center>
                  </Box>
                )}
              </Dropzone>
            </Stack>
            <HStack>
              <Stack spacing={-1} py="10px" w="full">
                <FormLabel
                  fontWeight={'600'}
                  fontSize={'lg'}
                >
                  Source Language
                </FormLabel>
                <Select value={sourceLanguage} onChange={(e) => setSourceLanguage((e.target.value as SupportedLanguage))}>
                  {SupportedLanguages.map((lang, idx) => (
                    <option key={idx} value={lang}>{lang}</option>
                  ))}
                </Select>
              </Stack>
              <Spacer />
              <Stack spacing={-1} py="10px" w="full">
                <FormLabel
                  fontWeight={'600'}
                  fontSize={'lg'}
                >
                  Target Language
                </FormLabel>
                <Select value={targetLanguage} onChange={(e) => setTargetLanguage((e.target.value as SupportedLanguage))}>
                  {SupportedLanguages.map((lang, idx) => (
                    <option key={idx} value={lang}>{lang}</option>
                  ))}
                </Select>
              </Stack>
            </HStack>
          </FormControl>
        </ModalBody>
        <ModalFooter alignSelf={"center"}>
          <HStack w="full">
            <Button
              leftIcon={<Check />}
              w="full"
              colorScheme="green"
              onClick={createProject}
              px="25px"
              isDisabled={!formValid}
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

export default NewProjectModal;
