import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useAuth } from "../hooks/useAuth";
import Chaos from "@/assets/doodle-stars.png";
import FileUpload from "./FileUpload"; // Import the FileUpload component
import axios from "axios";
import { useToast } from "@/components/ui/use-toast";

interface File {
  FileID: string;
  UserID: string;
  FileName: string;
  FileSize: number;
  FileType: string;
  CreatedAt: string;
  UpdatedAt: string;
}

const Dashboard = () => {
  const { user, signOut, getAuthToken } = useAuth();
  const [showFileManager, setShowFileManager] = useState(false);
  const [files, setFiles] = useState<File[]>([]);
  const { toast } = useToast();

  const fetchFiles = useCallback(async () => {
    try {
      const token = await getAuthToken();
      const response = await axios.get(
        "https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/chaosfiles-list-files",
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      setFiles(response.data);
    } catch (error) {
      console.error("Failed to fetch files:", error);
    }
  }, [getAuthToken]);

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  const handleUploadComplete = () => {
    fetchFiles();
  };

  const handleDownload = async (fileID: string) => {
    try {
      const token = await getAuthToken();
      const response = await axios.get(
        `https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/download-url?fileID=${fileID}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const { downloadUrl, fileName, contentType } = response.data;

      // tempo anchor element to trigger the download
      const link = document.createElement("a");
      link.href = downloadUrl;
      link.download = fileName;
      link.type = contentType;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error("Failed to download file: ", error);
    }
  };

  const handleDelete = async (fileID: string) => {
    try {
      const token = await getAuthToken();
      await axios.delete(
        `https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/chaosfiles-delete-file/${fileID}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      toast({
        title: "Success",
        description: "File deleted successfully.",
      });
      fetchFiles();
    } catch (error) {
      console.error("Failed to delete file: ", error);
      toast({
        title: "Error",
        description: "Failed to delete file. Please try again.",
        variant: "destructive",
      });
    }
  };

  const UserMenu = () => (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">Hi, {user?.username}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setShowFileManager(true)}>
          File Manager
        </DropdownMenuItem>
        <DropdownMenuItem onClick={signOut}>Sign out</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );

  const FileManager = () => (
    <div className="mt-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Files</h2>
        <div>
          <Input className="mr-2" placeholder="Search files, folders" />
          <Button variant="outline" className="mr-2">
            Create
          </Button>
          <Button>Upload here</Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Size</TableHead>
            <TableHead>Uploaded</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {files.map((file) => (
            <TableRow key={file.FileID}>
              <TableCell>{file.FileName}</TableCell>
              <TableCell>{file.FileSize}</TableCell>
              <TableCell>{new Date(file.CreatedAt).toLocaleString()}</TableCell>
              <TableCell>
                <FileActionsMenu fileID={file.FileID} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );

  const FileActionsMenu = ({ fileID }: { fileID: string }) => (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">...</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem>Details</DropdownMenuItem>
        <DropdownMenuItem onClick={() => handleDownload(fileID)}>
          Download
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => handleDelete(fileID)}>
          Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-8">
        <div className="flex items-center">
          <img src={Chaos} alt="ChaosFiles Logo" className="h-10 mr-4" />
          <h1 className="text-3xl font-bold">Welcome!</h1>
        </div>
        <UserMenu />
      </div>
      {showFileManager ? (
        <FileManager />
      ) : (
        <>
          <p className="text-gray-600 mb-4">
            Share your files with ChaosFiles.com. Upload as much data as you
            want, download without speed limits and share files with friends
            without additional costs!
          </p>
          <FileUpload OnUploadComplete={handleUploadComplete} />
        </>
      )}
    </div>
  );
};

export default Dashboard;
