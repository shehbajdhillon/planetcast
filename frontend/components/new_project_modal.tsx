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
  useColorModeValue,
  useDisclosure,
} from '@chakra-ui/react';
import { Check } from 'lucide-react';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import Dropzone from "react-dropzone";
import NProgress from 'nprogress';

import Image from 'next/image';

const CREATE_PROJECT = gql`
  mutation CreateProject($teamSlug: String!, $title: String!, $sourceLanguage: SupportedLanguage!, $sourceMedia: Upload!) {
    createProject(teamSlug: $teamSlug, title: $title, sourceLanguage: $sourceLanguage, sourceMedia: $sourceMedia) {
      id
      title
    }
  }
`;

interface NewProjectModalProps {
  refetch: () => void;
  teamSlug: string;
};

const MB = 1 << 20;

const NewProjectModal: React.FC<NewProjectModalProps> = (props) => {

  const { teamSlug, refetch } = props;

  const [title, setTitle] = useState("");
  const [sourceMedia, setSourceMedia] = useState<File>();
  const [sourceLanguage, setSourceLanguage] = useState<SupportedLanguage>("ENGLISH");

  const [formValid, setFormValid] = useState(false);

  const [createProjectMutation, { data, loading }] = useMutation(CREATE_PROJECT);

  const { onOpen, isOpen, onClose } = useDisclosure();

  const imgSrc = useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg');
  const borderColor = useColorModeValue('gray.300', 'whiteAlpha.500');
  const bgColor = useColorModeValue('white', 'black');

  const createProject = async () => {
    const res = await createProjectMutation({ variables: { title, teamSlug, sourceLanguage, sourceMedia } });
    if (res) {
      refetch();
      onClose();
    }
  };

  const router = useRouter();

  useEffect(() => {
    if (data?.createProject) {
      router.push(`/${teamSlug}/${data?.createProject.id}`)
    };
  }, [data, router, teamSlug]);

  useEffect(() => {
    if (loading) NProgress.start();
    else NProgress.done();
  }, [loading]);

  useEffect(() => {
    const checkFormValid = () => {
      if (!title.length) return false;
      if (!sourceMedia) return false;
      if (SupportedLanguages.indexOf(sourceLanguage) === -1) return false;
      return true;
    };
    setFormValid(checkFormValid());
  }, [title, sourceMedia, sourceLanguage]);

  useEffect(() => {
    setTitle("");
    setSourceMedia(undefined);
  }, [isOpen, onClose]);

  return (
    <Box>
      <Box
        p={5}
        onClick={onOpen}
        rounded={"lg"}
        _hover={{
          borderColor: borderColor,
          borderWidth: "1px",
          boxShadow: 'lg',
          bg: bgColor,
        }}
        cursor={"pointer"}
      >
        <Center h="full" flexDirection={"column"}>
          <Image
            src={imgSrc}
            width={70}
            height={100}
            style={{ borderRadius: "20px" }}
            alt='planet cast logo'
          />
          <Text>New Project</Text>
        </Center>
      </Box>
      <Modal isOpen={isOpen} onClose={onClose} isCentered size={{ base: "lg", md: "2xl" }}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Project</ModalHeader>
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
                  Media File (Max Size: 25 MB)
                </FormLabel>
                <Dropzone
                  accept={{ 'video/mp4': ['.mp4', '.MP4'] }}
                  minSize={0}
                  maxSize={25 * MB}
                  onDrop={(acceptedFiles) => setSourceMedia(acceptedFiles[0])}
                >
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
                isDisabled={!formValid || loading}
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
    </Box>
  );
};

export default NewProjectModal;
