import { Project, SupportedLanguage, SupportedLanguages, Transformation } from "@/types";
import { gql, useMutation } from "@apollo/client";
import {
  Box,
  Button,
  FormLabel,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Select,
  Stack,
  useDisclosure
} from "@chakra-ui/react";
import { PlusIcon } from "lucide-react";
import { useEffect, useState } from "react";

interface NewTransformationModelProps {
  project: Project;
  refetch: () => void;
};

const CREATE_TRANSLATION = gql`
  mutation CreateTranslation($projectId: Int64!, $targetLanguage: SupportedLanguage!) {
    createTranslation(projectId: $projectId, targetLanguage: $targetLanguage) {
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

  const createTranslation = async () => {
    const variables = { projectId: project.id, targetLanguage }
    const res = await createTranslationMutation({ variables });
    if (res) {
      refetch();
      onClose();
    }
    return res
  };

  useEffect(() => {
    setTargetLanguage(undubbedLanguages[0]);
  }, [isOpen]);

  return (
    <Box>
      <Button leftIcon={<PlusIcon />} onClick={onOpen} variant={"outline"}>New Dubbing</Button>
      <Modal isOpen={isOpen} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Dubbing</ModalHeader>
          <ModalCloseButton />
          <ModalBody overflow={"auto"}>
            <Stack w="full" h={"full"} direction={"row"}>
              <Box mx="auto">
                <Stack spacing={2}>
                  <FormLabel
                    fontWeight={'600'}
                    fontSize={'lg'}
                  >
                    Target Language
                  </FormLabel>
                  <Select value={targetLanguage} onChange={(e) => setTargetLanguage((e.target.value as SupportedLanguage))}>
                    {undubbedLanguages.map((lang, idx) => (
                      <option key={idx} value={lang}>{lang}</option>
                    ))}
                  </Select>
                  <Button hidden={!targetLanguage} onClick={createTranslation} isDisabled={loading}>
                    Submit
                  </Button>
                </Stack>
              </Box>
            </Stack>
          </ModalBody>
          <ModalFooter w="full">
            <Button colorScheme="gray" onClick={onClose}>Close</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};

export default NewTransformationModel;
