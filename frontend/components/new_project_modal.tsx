import { SupportedLanguage, SupportedLanguages, UploadOption } from '@/types';
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
  Checkbox,
  RadioGroup,
  Radio,
} from '@chakra-ui/react';
import { Check } from 'lucide-react';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import Dropzone from "react-dropzone";
import NProgress from 'nprogress';

import Image from 'next/image';
import { extractVideoID } from '@/utils';

const CREATE_PROJECT = gql`
  mutation CreateProject($teamSlug: String!, $title: String!, $sourceMedia: Upload, $uploadOption: UploadOption!, $youtubeLink: String, $initialLipSync: Boolean!, $initialTargetLanguage: SupportedLanguage) {
    createProject(teamSlug: $teamSlug, title: $title, sourceMedia: $sourceMedia, initialLipSync: $initialLipSync, initialTargetLanguage: $initialTargetLanguage, uploadOption: $uploadOption, youtubeLink: $youtubeLink) {
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
  const [youtubeLink, setYoutubeLink] = useState("");

  const [formValid, setFormValid] = useState(false);

  const [createProjectMutation, { data, loading }] = useMutation(CREATE_PROJECT);

  const { onOpen, isOpen, onClose } = useDisclosure();

  const [initialTargetLang, setInitialTargetLang] = useState<SupportedLanguage>("SPANISH");

  const [enableDubbing, setEnableDubbing] = useState(false);

  const [lipSync, setLipSync] = useState(false);

  const [uploadOption, setUploadOption] = useState<UploadOption>("FILE_UPLOAD");

  const imgSrc = useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg');
  const borderColor = useColorModeValue('blackAlpha.400', 'whiteAlpha.500');
  const bgColor = useColorModeValue('white', 'black');

  const createProject = async () => {
    const variables = {
      title,
      teamSlug,
      sourceMedia,
      initialLipSync: lipSync,
      initialTargetLanguage: enableDubbing ? initialTargetLang : undefined,
      uploadOption,
      youtubeLink,
    }
    const res = await createProjectMutation({
      variables
    });
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
      if (uploadOption === "FILE_UPLOAD" && !sourceMedia) return false;
      if (uploadOption === "YOUTUBE_LINK" && !extractVideoID(youtubeLink)) return false
      return true;
    };
    setFormValid(checkFormValid());
  }, [title, sourceMedia, uploadOption, youtubeLink]);

  useEffect(() => {
    setTitle("");
    setSourceMedia(undefined);
    setYoutubeLink("");
    setEnableDubbing(false);
    setLipSync(false);
  }, [isOpen, onClose]);

  useEffect(() => {
    setYoutubeLink("");
    setSourceMedia(undefined);
  }, [uploadOption]);

  return (
    <Box>
      <Box
        p={5}
        borderWidth={"1px"}
        h={{ base: "186px", md: "203px" }}
        w={{ base: "330px", md: "360px" }}
        onClick={onOpen}
        rounded={"lg"}
        _hover={{
          borderColor: borderColor,
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
                  Project Title
                </FormLabel>
                <Input
                  onChange={(e) => setTitle(e.target.value)}
                  value={title}
                  placeholder='PlanetCast Podcast Episode #3'
                />
              </Stack>

              <Stack spacing={-1} py="10px">
                <RadioGroup value={uploadOption} onChange={(val) => setUploadOption(val as UploadOption)}>
                  <HStack>
                    <Radio value="FILE_UPLOAD">
                      File Upload
                    </Radio>
                    <Radio value="YOUTUBE_LINK">
                      YouTube Link
                    </Radio>
                  </HStack>
                </RadioGroup>
              </Stack>

              { uploadOption === 'FILE_UPLOAD' &&
                <Stack spacing={-1} py="10px">
                  <FormLabel
                    fontWeight={'600'}
                    fontSize={'lg'}
                  >
                    Media File (Max Size: 500 MB)
                  </FormLabel>
                  <Dropzone
                    accept={{ 'video/mp4': ['.mp4', '.MP4'] }}
                    minSize={0}
                    maxSize={500 * MB}
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
              }

              { uploadOption === "YOUTUBE_LINK" &&
                <Stack spacing={-1} py="10px">
                  <FormLabel
                    fontWeight={'600'}
                    fontSize={'lg'}
                  >
                    YouTube Link
                  </FormLabel>
                  <Input
                    onChange={(e) => setYoutubeLink(e.target.value)}
                    value={youtubeLink}
                    placeholder='https://www.youtube.com/watch?v=dQw4w9WgXcQ'
                  />
                </Stack>
              }


            </FormControl>
            <HStack>
              <Stack spacing={2} py="10px" w="full">
                <FormLabel
                  fontWeight={'600'}
                  fontSize={'lg'}
                >
                  <HStack spacing={3}>
                    <Text>Dub video to another language after transcription</Text>
                    <Checkbox checked={enableDubbing} onChange={(e) => setEnableDubbing(e.target.checked)} />
                  </HStack>
                  <Text fontSize={'sm'} fontWeight={'light'} fontStyle={'italic'} hidden={enableDubbing}>
                    Check this box to have the video automatically dubbed into another language of your choice after transcription is complete.
                    You can also add dubbing languages later if needed.
                  </Text>
                </FormLabel>
                { enableDubbing &&
                  <>
                    <FormLabel
                      fontWeight={'600'}
                      fontSize={'lg'}
                    >
                      Dubbing Language
                    </FormLabel>
                    <Select value={initialTargetLang} onChange={(e) => setInitialTargetLang((e.target.value as SupportedLanguage))}>
                      {SupportedLanguages.map((lang, idx) => (
                        <option key={idx} value={lang}>{lang}</option>
                      ))}
                    </Select>
                    <Checkbox isChecked={lipSync} onChange={() => setLipSync(curr => !curr)}>
                      Enable Lip Syncing (Experimental)
                    </Checkbox>
                  </>
                }
              </Stack>
              <Spacer />
            </HStack>
          </ModalBody>
          <ModalFooter w="full">
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
