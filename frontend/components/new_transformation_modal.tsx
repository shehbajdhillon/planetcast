import { Project } from "@/types";
import {
  Box,
  Button,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure
} from "@chakra-ui/react";
import { PlusIcon } from "lucide-react";
import { useEffect } from "react";

interface NewTransformationModelProps {
  project: Project;
};

const NewTransformationModel: React.FC<NewTransformationModelProps> = (props) => {

  const { project } = props;
  const { onOpen, isOpen, onClose } = useDisclosure();

  const sourceTransformation = project?.transformations?.find(t => t.isSource);

  useEffect(() => {
    console.log({ sourceTransformation });
  }, [sourceTransformation]);

  return (

    <Box>
      <IconButton aria-label="add new dubbing" icon={<PlusIcon />} variant={"ghost"} onClick={onOpen} />
      <Modal isOpen={isOpen} onClose={onClose} size={"6xl"} isCentered>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Dubbing</ModalHeader>
          <ModalCloseButton />

          <ModalBody>
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
