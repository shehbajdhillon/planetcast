import { Project, Segment, SupportedLanguage, SupportedLanguages, Transformation } from "@/types";
import { formatTime } from "@/utils";
import {
  Box,
  Button,
  FormLabel,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Select,
  Stack,
  Text,
  VStack,
  useDisclosure
} from "@chakra-ui/react";
import { PlusIcon } from "lucide-react";
import { useEffect, useState } from "react";

interface NewTransformationModelProps {
  project: Project;
};

const NewTransformationModel: React.FC<NewTransformationModelProps> = (props) => {

  const { project } = props;
  const { onOpen, isOpen, onClose } = useDisclosure();

  const sourceTransformation: Transformation | undefined  = project?.transformations?.find(t => t.isSource);
  const sourceTranscript = sourceTransformation?.transcript && JSON.parse(sourceTransformation.transcript)
  const sourceSegments: Segment[] = sourceTranscript?.segments;
  const [targetLanguage, setTargetLanguage] = useState<SupportedLanguage>();

  useEffect(() => {
    setTargetLanguage(undefined);
  }, [isOpen]);

  return (
    <Box>
      <IconButton aria-label="add new dubbing" icon={<PlusIcon />} variant={"ghost"} onClick={onOpen} />
      <Modal isOpen={isOpen} onClose={onClose} size={"6xl"} isCentered>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Dubbing</ModalHeader>
          <ModalCloseButton />

          <ModalBody overflow={"auto"} maxH={"80vh"}>
            <Stack w="full" h={"full"} direction={"row"}>

              <VStack w="50%">
                <Button variant={"outline"} pointerEvents={"none"}>{sourceTransformation?.targetLanguage}</Button>
                <VStack>
                  {sourceSegments.map((segment: Segment, idx: number) => (
                    <Button
                      key={idx}
                      rounded="10px"
                      whiteSpace={'normal'}
                      height="auto"
                      blockSize={'auto'}
                      w="full"
                      justifyContent="left"
                      leftIcon={<Text>{formatTime(segment.start)} - {formatTime(segment.end)}</Text>}
                      variant={'outline'}
                    >
                      <Text
                        key={idx}
                        textAlign={"left"}
                        padding={2}
                      >
                        { segment.text.trim() }
                      </Text>
                    </Button>
                  ))}
                </VStack>
              </VStack>

              <Box mx="auto">
                <Stack spacing={2}>
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
                  <Button hidden={!targetLanguage}>
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
