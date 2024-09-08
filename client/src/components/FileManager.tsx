import React, { useCallback, useEffect, useState } from "react";
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
import axios from "axios";
import { useToast } from "@/components/ui/use-toast";
import Navbar from "./Navbar";

interface File {
  FileID: string;
  UserID: string;
  FileName: string;
  FileSize: number;
  FileType: string;
  CreatedAt: string;
  UpdatedAt: string;
}

const FileManager: React.FC = () => {
  const { getAuthToken } = useAuth();
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

  const handleDownload = async (fileID: string) => {
    try {
      const token = await getAuthToken();
      const response = await axios.get(
        `https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/download-url?fileID=${fileID}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const { downloadUrl } = response.data;

      window.open(downloadUrl, "_blank");
    } catch (error) {
      console.error("Failed to download file: ", error);
      toast({
        title: "Error",
        description: "Failed to download file. Please try again.",
        variant: "destructive",
      });
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
    <div className="mt-4">
      <Navbar />
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Files</h2>
        <div>
          <Input className="mr-2" placeholder="Search files" />
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
              <TableCell>
                {(file.FileSize / (1024 * 1024)).toFixed(2)} MB
              </TableCell>
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
};

export default FileManager;
