import { useRouter } from "next/router";
import { useState } from "react";
import SingleActionModal from "@/components/single_action_modal";
import { gql, useMutation } from "@apollo/client";
import { Heading, Box, Button, Grid, GridItem, HStack, Spacer, Stack, Text, VStack, useDisclosure } from "@chakra-ui/react";

const DELETE_PROJECT = gql`
  mutation DeleteProject($projectId: Int64!) {
    deleteProject(projectId: $projectId) {
      id
    }
  }
`;

interface ProjectSettingsTabProps {
  projectId: number;
  teamSlug: string;
  refetch: () => void;
};

const ProjectSettingsTab: React.FC<ProjectSettingsTabProps> = ({ refetch, projectId, teamSlug }) => {

  const [deleteProjectMutation, { loading }] = useMutation(DELETE_PROJECT);

  const router = useRouter();

  const deleteProject = async () => {
    const res = await deleteProjectMutation({ variables: { projectId } });
    if (res) {
      refetch();
      router.push(`/dashboard/${teamSlug}`);
    }
  };

  const { isOpen, onClose, onOpen } = useDisclosure();

  const [tabIdx, setTabIdx] = useState(0);

  const RenderTabButtons = () => (
    <VStack w="full" alignItems={"flex-start"} px="10px" spacing={"10px"}>
      <Button
        w="full"
        variant={"ghost"}
        onClick={() => setTabIdx(0)}
        borderWidth={tabIdx === 0 ? '1px' : ''}
        justifyContent={"left"}
      >
        General
      </Button>
    </VStack>
  )

  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1450px"}>
        <Grid
          templateAreas={{
            base: `
              "main"
            `,
            lg: `"sidebar main"`
          }}
          gridTemplateColumns={{ base: "1fr", lg: "1fr 4fr"}}
          w="full"
          h="full"
          gap="10px"
        >
          <GridItem area={"sidebar"} display={{ base: "none", lg: "block" }}>
            <RenderTabButtons />
          </GridItem>
          <GridItem area={"main"} w="full">
            {tabIdx === 0 &&
              <VStack alignItems={{ lg: "flex-start" }}>
                <Heading>Danger Zone</Heading>
                <Stack
                  direction={"column"}
                  w="full"
                  borderColor={"red.200"}
                  borderWidth={"1px"}
                  padding={"25px"}
                  rounded={"lg"}
                >
                  <HStack>
                    <Box>
                      <Text>Delete this Project</Text>
                      <Text>The original video and all the dubbings will be deleted. This action is not reversible.</Text>
                    </Box>
                    <Spacer />
                    <Box>
                      <SingleActionModal
                        heading={"Delete Project"}
                        action={() => deleteProject()}
                        loading={loading}
                        isOpen={isOpen}
                        onClose={onClose}
                      >
                        Are you sure you want to delete this Project? This will delete the original video and all the dubbings generated. This action is irreversible.
                      </SingleActionModal>
                      <Button colorScheme="red" onClick={onOpen}>
                        Delete Project
                      </Button>
                    </Box>
                  </HStack>
                </Stack>
              </VStack>
            }
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

export default ProjectSettingsTab;
