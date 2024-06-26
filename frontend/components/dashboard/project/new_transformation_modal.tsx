import { Project, SupportedLanguage, SupportedLanguages, Transformation } from "@/types";
import { gql, useMutation } from "@apollo/client";
import {
  Box,
  Button,
  Text,
  Checkbox,
  FormLabel,
  HStack,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Radio,
  RadioGroup,
  Select,
  Stack,
  useDisclosure,
  Center
} from "@chakra-ui/react";
import { PlusIcon } from "lucide-react";
import { useEffect, useState } from "react";

interface NewTransformationModelProps {
  project: Project;
  refetch: () => void;
};

const CREATE_TRANSLATION = gql`
  mutation CreateTranslation($projectId: Int64!, $targetLanguage: String!, $lipSync: Boolean!, $gender: String!) {
    createTranslation(projectId: $projectId, targetLanguage: $targetLanguage, lipSync: $lipSync, gender: $gender) {
      id
      projectId
    }
  }
`;

const NewTransformationModel: React.FC<NewTransformationModelProps> = (props) => {

  const { project, refetch } = props;
  const { onOpen, isOpen, onClose } = useDisclosure();

  const transformations: Transformation[] | undefined  = project?.transformations;
  const dubbedLanguages = transformations?.map((t) => t.targetLanguage);
  const undubbedLanguages = SupportedLanguages.filter((lang) => dubbedLanguages.indexOf(lang) == -1)

  const [targetLanguage, setTargetLanguage] = useState<SupportedLanguage>(undubbedLanguages[0]);

  const [createTranslationMutation, { loading }] = useMutation(CREATE_TRANSLATION);
  const [lipSync, setLipSync] = useState(false);
  const [gender, setGender] = useState("male");

  const createTranslation = async () => {
    const variables = { projectId: project.id, targetLanguage, lipSync, gender }
    const res = await createTranslationMutation({ variables });
    if (res) {
      refetch();
      onClose();
    }
    return res
  };

  useEffect(() => {
    setTargetLanguage(undubbedLanguages[0]);
    setLipSync(false);
    setGender("male");
  }, [isOpen]);

  return (
    <Box>
      <Button leftIcon={<PlusIcon />} onClick={onOpen} variant={"outline"}>New Dubbing</Button>
      <Modal isOpen={isOpen} onClose={onClose} isCentered size={"xl"}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Dubbing</ModalHeader>
          <ModalCloseButton />
          <ModalBody overflow={"auto"}>
            <Stack w="full" h={"full"} direction={"row"}>
              <Box mx="auto">
                <Stack spacing={5}>
                  <FormLabel
                    fontWeight={'600'}
                    fontSize={'lg'}
                  >
                    Dubbing Language
                  </FormLabel>
                  <Select value={targetLanguage} onChange={(e) => setTargetLanguage((e.target.value as SupportedLanguage))}>
                    {undubbedLanguages.map((lang, idx) => (
                      <option key={idx} value={lang}>{lang}</option>
                    ))}
                  </Select>
                  <Checkbox isChecked={lipSync} onChange={() => setLipSync(curr => !curr)}>
                    Enable Lip Syncing (Experimental)
                  </Checkbox>
                  <RadioGroup value={gender} onChange={setGender}>
                    <HStack>
                      <Radio value='male'>Male</Radio>
                      <Radio value='female'>Female</Radio>
                    </HStack>
                    <Text fontSize={'sm'} fontWeight={'light'} fontStyle={'italic'}>
                      We are currently only offering single speaker dubbing in male or female voices,
                      but are working on adding more voice options and multi-speaker dubbing capabilities.
                    </Text>
                  </RadioGroup>
                  <Text>Dubbing credits required: {project.dubbingCreditsRequired} minutes</Text>
                </Stack>
              </Box>
            </Stack>
          </ModalBody>
          <ModalFooter alignItems={"center"}>
            <Center w="full">
              <HStack>
                <Button
                  hidden={!targetLanguage}
                  onClick={createTranslation}
                  isDisabled={loading}
                  colorScheme="green"
                >
                  Submit
                </Button>
                <Button colorScheme="gray" onClick={onClose}>
                  Close
                </Button>
              </HStack>
            </Center>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};

export default NewTransformationModel;
