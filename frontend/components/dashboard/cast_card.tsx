import {
  Box,
  HStack,
  Spacer,
  Text,
  Button,
  useColorModeValue,
} from "@chakra-ui/react";
import { useRouter } from "next/router";

interface ProjectCardProps {
  projectId: number;
  teamSlug: string;
  title: string;
  status: "DRAFT" | "PROCESSING" | "DONE";
};

const ProjectCard: React.FC<ProjectCardProps> = (props) => {

  const { title, status, projectId, teamSlug } = props;

  const router = useRouter();

  return (
    <Box
      borderWidth={"1px"}
      p={6}
      maxW="400px"
      minW="270px"
      rounded={"lg"}
      _hover={{
        borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
        boxShadow: 'lg',
        bg: useColorModeValue('white', 'whiteAlpha.100'),
      }}
      cursor={"pointer"}
      onClick={() => router.push(`/dashboard/${teamSlug}/${projectId}`)}
    >
      <HStack>
        <Text
          textTransform="capitalize"
          maxW="400px"
          fontWeight={700}
          fontSize={'lg'}
          letterSpacing={1.1}
          noOfLines={1}
        >
          {title}
        </Text>
      </HStack>
      <Spacer />
      <HStack pt="20px">
        <Button
          borderWidth="1px"
          size={'xs'}
          textTransform="capitalize"
          fontWeight="medium"
          alignContent="right"
          pointerEvents={"none"}
        >
          {status}
        </Button>
      </HStack>
    </Box>
  );
};

export default ProjectCard;
